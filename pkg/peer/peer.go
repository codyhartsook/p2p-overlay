package peer

import (
	"context"
	"fmt"
	"p2p-overlay/pkg/cable"
	link_monitor "p2p-overlay/pkg/link-monitor"
	"p2p-overlay/pkg/pubsub"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	pb "p2p-overlay/pkg/grpc"

	"google.golang.org/grpc"
)

const (
	grpcPort           = 4224
	monitoringInterval = 10
)

type Peer struct {
	cable      cable.Cable
	grpcClient pb.PeersClient
	pubsub.Subscriber
	link_monitor.Monitor
	natsHost  string
	grpcAddr  string
	zone      string
	nodeZones map[string]string
}

func NewPeer(peerCableType, brokerHost string, hostZone string) *Peer {
	p := &Peer{zone: hostZone}

	p.natsHost = brokerHost
	p.grpcAddr = fmt.Sprintf("%s:%d", brokerHost, grpcPort)

	p.cable = cable.NewCable(peerCableType)

	err := p.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}

	p.connectToBroker()
	p.RegisterNatsSubscriber(p.natsHost)
	p.SubscribeToChannels(p.updateLocalPeers)

	return p
}

func (p *Peer) connectToBroker() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(p.grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewPeersClient(conn)
	p.grpcClient = client

	log.Print("connected to broker over grpc")
}

func (p *Peer) updateLocalPeers(peers []pubsub.PubPeer) {
	log.Printf("new peers broadcasted.")
	ctx := context.TODO()

	key := p.cable.GetPubKey()
	peerConfs := make([]wgtypes.PeerConfig, len(peers))

	member := false
	for i, peer := range peers {
		p.nodeZones[peer.Peer.AllowedIPs[0].IP.String()] = peer.Zone
		peerConfs[i] = peer.Peer
		if peer.Peer.PublicKey.String() == key {
			member = true
			break
		}
	}

	if !member {
		peerConfs = make([]wgtypes.PeerConfig, 0)
		log.Println("delete self signal received.")
	}

	p.cable.SyncPeers(ctx, peerConfs)
}

func (p *Peer) getNodeZone(ip string) string {
	return p.nodeZones[ip]
}

func (p *Peer) RegisterSelf() {
	// Perform config handshake with broker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	conf := p.cable.GetLocalConfig()
	req, err := p.cable.PeerConfigToProtobuf(conf)
	if err != nil {
		log.Printf("error converting peer config to protobuf: %v", err)
	}

	req.Zone = p.zone
	brokerRes, err := p.grpcClient.RegisterPeer(ctx, req)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	if !brokerRes.Success {
		log.Fatalf("broker rejected peer registration")
	}

	log.Printf("overlay address provisioned: %s", brokerRes.Address)

	p.cable.SetAddress(brokerRes.Address)
	p.cable.AddrAdd()

	// monitor tunnel performance
	p.InitializeMonitoring(p.natsHost, brokerRes.Address, "peer")
	p.StartMonitor(monitoringInterval, p.cable.GetPeerTopology, p.getNodeZone)
}

func (p *Peer) unRegisterSelf() {
	// send grpc request to broker
	// remove local interfaces
}
