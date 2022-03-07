package broker

import (
	"context"
	"net"
	"p2p-overlay/pkg/cable"
	"sync"

	log "github.com/sirupsen/logrus"

	pb "p2p-overlay/pkg/grpc"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

const (
	natsHost = "localhost:4222"
	grpcAddr = "0.0.0.0:4224"
)

type Broker struct {
	cable     cable.Cable
	publisher *nats.EncodedConn
	pb.UnimplementedPeersServer
	mutex *sync.RWMutex
}

func NewBroker() *Broker {
	b := &Broker{}

	// start tunnel cable
	/*b.cable = cable.NewCable()

	err := b.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}*/

	return b
}

func (b *Broker) StartListeners() {
	// start nats
	b.registerPublisher()

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

func (b *Broker) registerPublisher() {
	nc, err := nats.Connect(natsHost)
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}

	conn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}

	b.publisher = conn
	log.Println("connected to nats")
}

// Implement the gRPC interface
func (b *Broker) RegisterPeer(ctx context.Context, peer *pb.RegisterPeersRequest) (*pb.RegisterPeersResponse, error) {
	log.Printf("Registering peer: %v", peer)
	return &pb.RegisterPeersResponse{}, nil
}

func (b *Broker) UnregisterPeer(ctx context.Context, peer *pb.UnregisterPeersRequest) (*pb.UnregisterPeersResponse, error) {
	log.Printf("Unregistering peer: %v", peer)
	return &pb.UnregisterPeersResponse{}, nil
}

/*
Listen on grpc for new peer requesets

Recieve peer config and register new peer

Return self config to new peer

Publish new peer config to peer-subscriber channel

*/
