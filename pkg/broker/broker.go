package broker

import (
	"context"
	"net"
	"p2p-overlay/pkg/cable"
	"p2p-overlay/pkg/pubsub"
	"sync"

	log "github.com/sirupsen/logrus"

	pb "p2p-overlay/pkg/grpc"

	"google.golang.org/grpc"
)

const (
	grpcAddr = "0.0.0.0:4224"
)

type Broker struct {
	cable cable.Cable
	pubsub.Publisher
	pb.UnimplementedPeersServer
	mutex *sync.RWMutex
}

func NewBroker(peerCableType string) *Broker {
	b := &Broker{mutex: &sync.RWMutex{}}

	// start tunnel cable
	b.cable = cable.NewCable(peerCableType)

	err := b.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}

	return b
}

func (b *Broker) StartListeners() {
	// start nats
	b.RegisterPublisher()

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
func (b *Broker) RegisterPeer(ctx context.Context, peer *pb.RegisterPeerRequest) (*pb.RegisterPeerResponse, error) {
	log.Info("registering peer")

	err := b.addPeerToLocalConfig(peer)
	if err != nil {
		log.Info("error adding peer to local config: %v", err)
		return &pb.RegisterPeerResponse{}, err
	}

	// publish to nats
	ctx = context.TODO()
	peers, err := b.cable.GetPeers(ctx)
	if err != nil {
		log.Info("error getting peers from local config: %v", err)
		return &pb.RegisterPeerResponse{}, err
	}

	b.BroadcastPeers(peers)
	return &pb.RegisterPeerResponse{}, nil
}

func (b *Broker) UnregisterPeer(ctx context.Context, peer *pb.UnregisterPeerRequest) (*pb.UnregisterPeerResponse, error) {
	log.Printf("Unregistering peer: %v", peer)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	err := b.removePeerFromLocalConfig(peer)
	if err != nil {
		log.Printf("error removing peer from local config: %v", err)
		return &pb.UnregisterPeerResponse{}, err
	}

	// publish to nats
	ctx = context.TODO()
	peers, err := b.cable.GetPeers(ctx)
	if err != nil {
		log.Info("error getting peers from local config: %v", err)
		return &pb.UnregisterPeerResponse{}, err
	}

	b.BroadcastPeers(peers)
	return &pb.UnregisterPeerResponse{}, nil
}

func (b *Broker) addPeerToLocalConfig(peer *pb.RegisterPeerRequest) error {
	log.Printf("adding peer %s to local config", peer)
	ctx := context.TODO()
	conf := cable.ProtobufPeerToConfig(peer.Peer)

	return b.cable.RegisterPeer(ctx, conf)
}

func (b *Broker) removePeerFromLocalConfig(peer *pb.UnregisterPeerRequest) error {
	log.Printf("removing peer %s from local config", peer)
	ctx := context.TODO()
	return b.cable.DeletePeer(ctx, peer.PublicKey)
}
