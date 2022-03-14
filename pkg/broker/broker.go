package broker

import (
	"context"
	"net"
	"p2p-overlay/pkg/cable"
	link_monitor "p2p-overlay/pkg/link-monitor"
	"p2p-overlay/pkg/pubsub"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	pb "p2p-overlay/pkg/grpc"
	addresses "p2p-overlay/pkg/subnet"

	"google.golang.org/grpc"
)

const (
	grpcAddr   = "0.0.0.0:4224"
	arangoHost = "127.0.0.1"
)

type Broker struct {
	cable cable.Cable
	pubsub.Publisher
	addresses.AddressDistribution
	link_monitor.Monitor
	pb.UnimplementedPeersServer
	mutex    *sync.RWMutex
	natsHost string
	nodes    map[string]pubsub.NodeSpec
	zone     string
}

func NewBroker(peerCableType string, brokerHost string, hostZone string) *Broker {
	b := &Broker{mutex: &sync.RWMutex{}, natsHost: brokerHost, zone: hostZone}

	b.InitializeAddresses()
	brokerAddr := b.GetBrokerAddress().String()

	// start wg tunnel agent
	b.cable = cable.NewCable(peerCableType)
	b.cable.SetAddress(brokerAddr)

	err := b.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}

	b.nodes = make(map[string]pubsub.NodeSpec)
	b.nodes[b.cable.GetPubKey()] = pubsub.NodeSpec{Address: brokerAddr, Zone: b.zone}

	// add address to wireguard tunnel
	b.cable.AddrAdd()

	// monitor tunnel performance
	b.InitializeMonitoring(arangoHost, brokerAddr, "broker")
	b.StartMonitor(10, b.getPeerSubnetAddrs, b.getNodeZone)

	return b
}

func (b *Broker) getNodeZone(ip string) string {
	for _, node := range b.nodes {
		if node.Address == ip {
			return node.Zone
		}
	}
	return ""
}

func (b *Broker) getPeerSubnetAddrs() []net.IP {
	ips := make([]net.IP, len(b.nodes))
	i := 0
	for _, node := range b.nodes {
		ips[i] = net.ParseIP(node.Address)
		i++
	}

	return ips
}

func (b *Broker) StartListeners() {
	// start nats publisher
	b.RegisterPublisher(b.natsHost)

	// start grpc server
	b.registerGrpc()
}

func (b *Broker) registerGrpc() {
	// start grpc server
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	pb.RegisterPeersServer(gs, b)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("grpc server started")
}

// Implements the gRPC inerface
func (b *Broker) RegisterPeer(ctx context.Context, peer *pb.Peer) (*pb.RegisterPeerResponse, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	address, err := b.GetAvailableAddress()
	if err != nil {
		log.Printf("error getting available address: %v", err)
		return &pb.RegisterPeerResponse{Success: false}, err
	}

	peer.Address = address.String()
	err = b.addPeerToLocalConfig(peer)
	if err != nil {
		log.Info("error adding peer to local config: %v", err)
		return &pb.RegisterPeerResponse{Success: false}, err
	}

	b.nodes[peer.PublicKey] = pubsub.NodeSpec{Address: peer.Address, Zone: peer.Zone}

	// broadcast local peers config to nats
	peers, err := b.cable.GetPeers(nil)
	if err != nil {
		log.Info("error getting peers from local config: %v", err)
		return &pb.RegisterPeerResponse{Success: false}, err
	}

	b.newBroadcast(peers)

	return &pb.RegisterPeerResponse{Success: true, Address: peer.Address}, nil
}

func (b *Broker) UnregisterPeer(ctx context.Context, peer *pb.UnregisterPeerRequest) (*pb.UnregisterPeerResponse, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	innerCtx := context.TODO()
	err := b.cable.DeletePeer(innerCtx, peer.PublicKey)
	if err != nil {
		log.Printf("error removing peer from local config: %v", err)
		return &pb.UnregisterPeerResponse{Success: false}, err
	}

	delete(b.nodes, peer.PublicKey)

	// broadcast updated peers list, peers will sync with new config thus removing this peer
	ctx = context.TODO()
	peers, err := b.cable.GetPeers(ctx)
	if err != nil {
		log.Info("error getting peers from local config: %v", err)
		return &pb.UnregisterPeerResponse{Success: false}, err
	}

	b.newBroadcast(peers)
	return &pb.UnregisterPeerResponse{Success: true}, nil
}

func (b *Broker) addPeerToLocalConfig(peer *pb.Peer) error {
	log.Printf("adding peer %s to local config", peer.Endpoint)

	ctx := context.TODO()
	conf, err := b.cable.ProtobufToPeerConfig(peer)
	if err != nil {
		log.Printf("error converting protobuf peer to config: %v", err)
		return err
	}

	return b.cable.RegisterPeer(ctx, conf)
}

func (b *Broker) newBroadcast(peers []wgtypes.PeerConfig) {
	// include broker config in broadcasted peers list
	myConf := b.cable.GetLocalConfig()
	peers = append(peers, myConf)

	peersList := make([]pubsub.PubPeer, len(peers))
	for i, peer := range peers {
		peersList[i] = pubsub.PubPeer{
			Peer:     peer,
			Metadata: pubsub.NodeSpec{Address: b.nodes[peer.PublicKey.String()].Address, Zone: b.nodes[peer.PublicKey.String()].Zone},
		}
	}

	b.BroadcastPeers(peersList)
}
