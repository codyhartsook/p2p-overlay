package broker

import (
	"context"
	"net"
	"p2p-overlay/pkg/cable"
	link_monitor "p2p-overlay/pkg/link-monitor"
	"p2p-overlay/pkg/pubsub"
	"sync"

	log "github.com/sirupsen/logrus"

	"p2p-overlay/pkg/addresses"
	pb "p2p-overlay/pkg/grpc"

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
	peers    map[string]net.IP
}

func NewBroker(peerCableType string, brokerHost string) *Broker {
	b := &Broker{mutex: &sync.RWMutex{}, natsHost: brokerHost}

	b.InitializeAddresses()
	brokerAddr := b.GetBrokerAddress().String()

	b.peers = make(map[string]net.IP)

	// start tunnel cable
	b.cable = cable.NewCable(peerCableType)
	b.cable.SetAddress(brokerAddr)

	err := b.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}

	// add address to wireguard tunnel
	b.cable.AddrAdd()

	// monitor tunnel performance
	b.InitializeMonitoring(arangoHost, brokerAddr, "broker")
	b.StartMonitor(10, b.cable.GetPeerTopology)

	return b
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

// Implement the gRPC inerface
func (b *Broker) RegisterPeer(ctx context.Context, peer *pb.Peer) (*pb.RegisterPeerResponse, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var err error
	address, registered := b.peers[peer.PublicKey]
	if !registered {
		address, err = b.GetAvailableAddress()
		if err != nil {
			log.Printf("error getting available address: %v", err)
			return &pb.RegisterPeerResponse{Success: false}, err
		}
		b.peers[peer.PublicKey] = address
	}

	peer.Address = address.String()
	err = b.addPeerToLocalConfig(peer)
	if err != nil {
		log.Info("error adding peer to local config: %v", err)
		return &pb.RegisterPeerResponse{Success: false}, err
	}

	// broadcast local peers config to nats
	ctx = context.TODO()
	peers, err := b.cable.GetPeers(ctx)
	if err != nil {
		log.Info("error getting peers from local config: %v", err)
		return &pb.RegisterPeerResponse{Success: false}, err
	}

	myConf := b.cable.GetLocalConfig()
	peers = append(peers, myConf)

	b.BroadcastPeers(peers)
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

	delete(b.peers, peer.PublicKey)

	// publish to nats
	ctx = context.TODO()
	peers, err := b.cable.GetPeers(ctx)
	if err != nil {
		log.Info("error getting peers from local config: %v", err)
		return &pb.UnregisterPeerResponse{Success: false}, err
	}

	b.BroadcastPeers(peers)
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
