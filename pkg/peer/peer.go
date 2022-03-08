package peer

import (
	"context"
	"p2p-overlay/pkg/cable"
	"p2p-overlay/pkg/pubsub"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	pb "p2p-overlay/pkg/grpc"

	"google.golang.org/grpc"
)

const (
	grpcAddr = "localhost:4224"
	natsAddr = "localhost:4222"
)

type Peer struct {
	cable      cable.Cable
	grpcClient pb.PeersClient
	pubsub.Subscriber
}

func NewPeer(peerCableType string) *Peer {
	p := &Peer{}

	p.cable = cable.NewCable(peerCableType)

	err := p.cable.Init()
	if err != nil {
		log.Fatalf("error initializing wireguard device: %v", err)
	}

	p.connectToBroker()
	p.RegisterNatsSubscriber()

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

	log.Print("connected to broker over grpc")
}

func (p *Peer) SubscribeToOverlayUpdates() {
	p.SubscribeToChannels(p.updateLocalPeers)
}

func (p *Peer) updateLocalPeers(peers []wgtypes.Peer) {
	log.Printf("updating peers in local config")
	ctx := context.TODO()

	p.cable.SyncPeers(ctx, peers)
}

func (p *Peer) RegisterSelf() {
	// Perform config handshake with broker
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	req := pb.RegisterPeerRequest{}
	brokerRes, err := p.grpcClient.RegisterPeer(ctx, &req)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	ctx = context.TODO()
	brokerConf := cable.ProtobufPeerToConfig(brokerRes.Peer)

	p.cable.RegisterPeer(ctx, brokerRes.DeviceName, brokerConf)
}

func (p *Peer) UnRegisterSelf() {

}
