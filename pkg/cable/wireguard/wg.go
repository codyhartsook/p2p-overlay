package wireguard

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"p2p-overlay/pkg/endpoint"
	pb "p2p-overlay/pkg/grpc"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	port = 4500

	DefaultDeviceName = "wg-overlay"
	PublicKey         = "publicKey"

	// KeepAliveInterval to use for wg peers.
	KeepAliveInterval = 10 * time.Second
)

type specification struct {
	PSK      string `default:"default psk"`
	NATTPort int    `default:"4500"`
}

type WGCtrl struct {
	client *wgctrl.Client
	link   netlink.Link
	mutex  sync.Mutex
	keys   map[string]string
	psk    *wgtypes.Key
	spec   *specification
}

func NewWGCtrl() (*WGCtrl, error) {
	w := WGCtrl{spec: new(specification)}

	if err := w.addLink(); err != nil {
		return nil, errors.Wrap(err, "failed to setup WireGuard link")
	}

	// Create the controller
	var err error
	if w.client, err = wgctrl.New(); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("wgctrl is not available on this system")
		}

		return nil, errors.Wrap(err, "failed to open wgctl client")
	}

	log.Info("WireGuard client created")

	var priv, pub, psk wgtypes.Key
	if psk, err = genPsk(w.spec.PSK); err != nil {
		return nil, errors.Wrap(err, "error generating pre-shared key")
	}

	w.psk = &psk

	if priv, err = wgtypes.GeneratePrivateKey(); err != nil {
		return nil, errors.Wrap(err, "error generating private key")
	}

	pub = priv.PublicKey()
	w.keys = make(map[string]string)
	w.keys[PublicKey] = pub.String()

	portInt := int(port)

	// Configure the device - still not up.
	peerConfigs := make([]wgtypes.PeerConfig, 0)
	cfg := wgtypes.Config{
		PrivateKey:   &priv,
		ListenPort:   &portInt,
		FirewallMark: nil,
		ReplacePeers: true,
		Peers:        peerConfigs,
	}

	if err = w.client.ConfigureDevice(DefaultDeviceName, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to configure WireGuard device")
	}

	log.Infof("Created local wg device %s with publicKey %s", DefaultDeviceName, pub)

	return &w, nil
}

func (w *WGCtrl) Init() error {
	log.Infof("Initializing wg device")

	_, err := net.InterfaceByName(DefaultDeviceName)
	if err != nil {
		return errors.Wrapf(err, "cannot get wireguard link by name %s", DefaultDeviceName)
	}

	d, err := w.client.Device(DefaultDeviceName)
	if err != nil {
		return errors.Wrap(err, "wgctrl cannot find WireGuard device")
	}

	k, err := keyFromMap(w.keys)
	if err != nil {
		return errors.Wrapf(err, "endpoint is missing public key %s", d.PublicKey)
	}

	if k.String() != d.PublicKey.String() {
		return fmt.Errorf("endpoint public key %s is different from device key %s", k, d.PublicKey)
	}

	// IP link set $DefaultDeviceName up.
	if err := netlink.LinkSetUp(w.link); err != nil {
		return errors.Wrap(err, "failed to bring up WireGuard device")
	}

	log.Infof("WireGuard device %s is up",
		w.link.Attrs().Name)

	return nil
}

func (w *WGCtrl) RegisterPeer(ctx context.Context, peer wgtypes.PeerConfig) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	peers := []wgtypes.PeerConfig{peer}
	err := w.client.ConfigureDevice(DefaultDeviceName, wgtypes.Config{
		ReplacePeers: false,
		Peers:        peers,
	})
	if err != nil {
		return err
	}

	// verify peer was added
	if p, err := w.peerByKey(&peer.PublicKey); err != nil {
		log.Errorf("Failed to verify peer configuration: %v", err)
	} else {
		// TODO verify configuration
		log.Infof("Peer configured, PubKey:%s, EndPoint:%s, AllowedIPs:%v", p.PublicKey, p.Endpoint, p.AllowedIPs)
	}

	return nil
}

func (w *WGCtrl) GetPeers(ctx context.Context) ([]wgtypes.PeerConfig, error) {
	return nil, nil
}

func (w *WGCtrl) GetLocalConfig() wgtypes.PeerConfig {
	ip := endpoint.GetLocalIP()

	dev, _ := w.client.Device(DefaultDeviceName)
	conf := wgtypes.PeerConfig{
		PublicKey:    dev.PublicKey,
		Endpoint:     &net.UDPAddr{IP: ip, Port: port},
		PresharedKey: w.psk,
	}

	log.Info("Local peer config: ", conf)

	return conf
}

func (w *WGCtrl) SyncPeers(ctx context.Context, peers []wgtypes.PeerConfig) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	dev, _ := w.client.Device(DefaultDeviceName)

	filteredPeers := make([]wgtypes.PeerConfig, 0)
	for _, peer := range peers {
		if dev.PublicKey == peer.PublicKey {
			continue
		}
		filteredPeers = append(filteredPeers, peer)
	}

	err := w.client.ConfigureDevice(DefaultDeviceName, wgtypes.Config{
		ReplacePeers: true,
		Peers:        filteredPeers,
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *WGCtrl) DeletePeer(ctx context.Context, publicKey string) error {
	key, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return errors.Wrapf(err, "failed to parse public key %s", publicKey)
	}

	peerCfg := []wgtypes.PeerConfig{
		{
			PublicKey: key,
			Remove:    true,
		},
	}

	err = w.client.ConfigureDevice(DefaultDeviceName, wgtypes.Config{
		ReplacePeers: false,
		Peers:        peerCfg,
	})

	if err != nil {
		return errors.Wrapf(err, "failed to remove WireGuard peer with key %s", key)
	}

	log.Infof("Removed WireGuard peer with key %s", key)
	return nil
}

func (w *WGCtrl) peerByKey(key *wgtypes.Key) (*wgtypes.Peer, error) {
	d, err := w.client.Device(DefaultDeviceName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find device %s", DefaultDeviceName)
	}

	for i := range d.Peers {
		if d.Peers[i].PublicKey.String() == key.String() {
			return &d.Peers[i], nil
		}
	}

	return nil, fmt.Errorf("peer not found for key %s", key)
}

func (w *WGCtrl) PeerConfigToProtobuf(conf wgtypes.PeerConfig) (*pb.Peer, error) {
	allowedIps := make([]string, len(conf.AllowedIPs))
	for i, ip := range conf.AllowedIPs {
		allowedIps[i] = ip.String()
	}
	peer := &pb.Peer{
		PublicKey:    conf.PublicKey.String(),
		Endpoint:     conf.Endpoint.String(),
		PresharedKey: conf.PresharedKey.String(),
		AllowedIps:   allowedIps,
	}

	return peer, nil
}

func (w *WGCtrl) ProtobufToPeerConfig(peer *pb.Peer) (wgtypes.PeerConfig, error) {
	log.Printf("converting protobuf peer to wireguard peer config: %v", peer)

	key, err := wgtypes.ParseKey(peer.PublicKey)
	if err != nil {
		return wgtypes.PeerConfig{}, errors.Wrapf(err, "failed to parse public key %s", key)
	}

	remoteIP := net.ParseIP(strings.Split(peer.Endpoint, ":")[0])
	if remoteIP == nil {
		return wgtypes.PeerConfig{}, errors.Wrapf(err, "failed to parse ip %s", remoteIP)
	}
	remotePort, err := strconv.Atoi(strings.Split(peer.Endpoint, ":")[1])
	if err != nil {
		return wgtypes.PeerConfig{}, errors.Wrapf(err, "failed to parse port %s", port)
	}

	ka := KeepAliveInterval
	allowedIps := parseSubnets(peer.AllowedIps)
	pc := wgtypes.PeerConfig{
		PublicKey:    key,
		Remove:       false,
		UpdateOnly:   false,
		PresharedKey: w.psk,
		Endpoint: &net.UDPAddr{
			IP:   remoteIP,
			Port: remotePort,
		},
		PersistentKeepaliveInterval: &ka,
		AllowedIPs:                  allowedIps,
		ReplaceAllowedIPs:           false,
	}

	return pc, nil
}

func parseSubnets(subnets []string) []net.IPNet {
	nets := make([]net.IPNet, 0, len(subnets))

	for _, sn := range subnets {
		_, cidr, err := net.ParseCIDR(sn)
		if err != nil {
			// This should not happen. Log and continue.
			log.Errorf("failed to parse subnet %s: %v", sn, err)
			continue
		}

		nets = append(nets, *cidr)
	}

	return nets
}

func keyFromMap(keys map[string]string) (*wgtypes.Key, error) {
	s, found := keys[PublicKey]
	if !found {
		return nil, fmt.Errorf("missing public key")
	}

	key, err := wgtypes.ParseKey(s)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse public key %s", s)
	}

	return &key, nil
}

func genPsk(psk string) (wgtypes.Key, error) {
	// Convert spec PSK string to right length byte array, using sha256.Size == wgtypes.KeyLen.
	pskBytes := sha256.Sum256([]byte(psk))
	return wgtypes.NewKey(pskBytes[:]) // nolint:wrapcheck // Let the caller wrap it
}

func (w *WGCtrl) addLink() error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = w.setLinuxWGLink()
	case "darwin":
		err = w.setDarwinWGLink()
	default:
		log.Fatalf("unsupported OS: %s", runtime.GOOS)
	}

	return err
}

func (w *WGCtrl) setDarwinWGLink() error {
	err := exec.Command("ip", "link", "set", DefaultDeviceName, "type", "wireguard").Run()
	if err != nil {
		return err
	}
	return nil
}

func (w *WGCtrl) setLinuxWGLink() error {
	// delete existing wg device if needed
	if link, err := netlink.LinkByName(DefaultDeviceName); err == nil {
		// delete existing device
		if err := netlink.LinkDel(link); err != nil {
			return errors.Wrap(err, "failed to delete existing WireGuard device")
		}
	}

	// Create the wg device (ip link add dev $DefaultDeviceName type wireguard).
	la := netlink.NewLinkAttrs()
	la.Name = DefaultDeviceName
	link := &netlink.GenericLink{
		LinkAttrs: la,
		LinkType:  "wireguard",
	}

	if err := netlink.LinkAdd(link); err == nil {
		w.link = link
	} else {
		return errors.Wrap(err, "failed to add WireGuard device")
	}

	return nil
}
