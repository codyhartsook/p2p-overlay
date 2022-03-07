package peer

import (
	"context"
	"log"
	"p2p-overlay/pkg/cable"
	"time"

	pb "p2p-overlay/pkg/grpc"

	"google.golang.org/grpc"
)

const (
	grpcAddr = "localhost:4224"
)

type Peer struct {
	cable      cable.Cable
	grpcClient pb.PeersClient
}

func NewPeer() *Peer {
	p := &Peer{}

	p.cable = cable.NewCable()

	err := p.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}

	p.connectToBroker()

	return p
}

func (p *Peer) connectToBroker() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewPeersClient(conn)
	p.grpcClient = client

	log.Print("connected...")
}

func (p *Peer) RegisterSelf() {
	//var localConfig wgtypes.PeerConfig
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	req := pb.RegisterPeersRequest{}
	res, err := p.grpcClient.RegisterPeer(ctx, &req)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	log.Printf("registered: %s", res)
}

func (p *Peer) UnRegisterSelf() {

}

/*
Send grpc request to broker with peer config

Recieve broker peer config and nats peer-subscriber channel

Register peer with broker config

Listen to peer-subscriber channel

*/
