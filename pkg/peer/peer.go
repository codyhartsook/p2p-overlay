package peer

import (
	"context"
	"fmt"
	"p2p-overlay/pkg/cable"
	"p2p-overlay/pkg/pubsub"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	pb "p2p-overlay/pkg/grpc"

	"google.golang.org/grpc"
)

const (
	grpcPort = 4224
)

type Peer struct {
	cable      cable.Cable
	grpcClient pb.PeersClient
	pubsub.Subscriber
	natsHost string
	grpcAddr string
}

func NewPeer(peerCableType, brokerHost string) *Peer {
	p := &Peer{}

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

func (p *Peer) updateLocalPeers(peers []wgtypes.PeerConfig) {
	log.Printf("updating peers in local config")
	ctx := context.TODO()

	p.cable.SyncPeers(ctx, peers)
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

	brokerRes, err := p.grpcClient.RegisterPeer(ctx, req)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	if !brokerRes.Success {
		log.Fatalf("broker rejected peer registration")
	}
}

func (p *Peer) UnRegisterSelf() {
	// send grpc request to broker
	// remove local interfaces
}
