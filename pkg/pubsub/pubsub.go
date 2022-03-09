package pubsub

import (
	"fmt"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	natsPort  = 4222
	PeerNodes = "network-peer-nodes"
)

type Publisher struct {
	conn *nats.EncodedConn
}

func (p *Publisher) RegisterPublisher(natsHost string) {
	natsAddr := fmt.Sprintf("%s:%d", natsHost, natsPort)
	nc, err := nats.Connect(natsAddr)
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}

	conn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}

	p.conn = conn
	log.Info("connected to nats server.")
}

func (p *Publisher) BroadcastPeers(peers []wgtypes.PeerConfig) {
	err := p.conn.Publish(PeerNodes, peers)
	if err != nil {
		log.Fatalf("error publishing to channel %s: %v", PeerNodes, err)
	}

	log.Info("published to channel %s", PeerNodes)
}

type Subscriber struct {
	subConn *nats.EncodedConn
}

func (s *Subscriber) RegisterNatsSubscriber(natsHost string) {
	natsAddr := fmt.Sprintf("%s:%d", natsHost, natsPort)
	nc, err := nats.Connect(natsAddr)
	if err != nil {
		log.Fatalf("error connecting to nats server: %v", err)
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("error creating encoder connection: %v", err)
	}

	log.Info("connected to nats server.")
	s.subConn = c
}

func (s *Subscriber) SubscribeToChannels(handler func([]wgtypes.PeerConfig)) {
	// Nats Async Ephemeral Consumer

	log.Printf("subscribing to channel %s", PeerNodes)
	_, err := s.subConn.Subscribe(PeerNodes, func(m []wgtypes.PeerConfig) {
		handler(m)
	})
	if err != nil {
		log.Fatalf("error subscribing to channel %s: %v", PeerNodes, err)
	}

	s.subConn.Flush()

	if err := s.subConn.LastError(); err != nil {
		log.Fatal(err)
	}
}
