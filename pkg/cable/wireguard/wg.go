package wireguard

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	port              = 4500
	DefaultDeviceName = "wg-overlay"
	PublicKey         = "publicKey"
)

type WGCtrl struct {
	client wgctrl.Client
	link   netlink.Link
	mutex  sync.Mutex
	keys   map[string]string
}

func NewWGCtrl() (*WGCtrl, error) {
	w := WGCtrl{}

	if err := w.addLink(); err != nil {
		return nil, errors.Wrap(err, "failed to setup WireGuard link")
	}

	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	w.client = *client
	w.keys = make(map[string]string)

	log.Info("WireGuard client created")

	var priv, pub wgtypes.Key
	if priv, err = wgtypes.GeneratePrivateKey(); err != nil {
		return nil, errors.Wrap(err, "error generating private key")
	}

	pub = priv.PublicKey()
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

	log.Info("Configuring WireGuard device")
	if err = w.client.ConfigureDevice(DefaultDeviceName, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to configure WireGuard device")
	}

	log.Infof("Created WireGuard %s with publicKey %s", DefaultDeviceName, pub)

	return &w, nil
}

func (w *WGCtrl) Init() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log.Infof("Initializing WireGuard device")

	l, err := net.InterfaceByName(DefaultDeviceName)
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

	log.Infof("WireGuard device %s, is up on i/f number %d, listening on port :%d, with key %s",
		w.link.Attrs().Name, l.Index, d.ListenPort, d.PublicKey)

	return nil
}

func (w *WGCtrl) RegisterPeers(ctx context.Context, deviceName string, peers []wgtypes.PeerConfig) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	err := w.client.ConfigureDevice(deviceName, wgtypes.Config{
		ReplacePeers: false,
		Peers:        peers,
	})
	if err != nil {
		return err
	}

	for _, peer := range peers {
		// verify peer was added
		if p, err := w.peerByKey(&peer.PublicKey); err != nil {
			log.Errorf("Failed to verify peer configuration: %v", err)
		} else {
			// TODO verify configuration
			log.Infof("Peer configured, PubKey:%s, EndPoint:%s, AllowedIPs:%v", p.PublicKey, p.Endpoint, p.AllowedIPs)
		}
	}

	return nil
}

func (w *WGCtrl) DeletePeers(ctx context.Context, deviceName string, publicKeys []string) error {
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

func (w *WGCtrl) addLink() error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = w.setWGLink()
	default:
		log.Fatalf("unsupported OS: %s", runtime.GOOS)
	}

	return err
}

func (w *WGCtrl) setWGLink() error {
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
