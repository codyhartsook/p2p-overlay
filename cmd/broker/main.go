package main

import (
	"flag"
	"p2p-overlay/pkg/broker"
)

func main() {

	cableType := flag.String("cable", "wireguard", "peer connection cable type")
	natsHost := flag.String("nats", "localhost", "nats host")

	flag.Parse()

	b := broker.NewBroker(*cableType, *natsHost)

	b.StartListeners()
}
