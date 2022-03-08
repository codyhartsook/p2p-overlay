package main

import (
	"flag"
	"p2p-overlay/pkg/broker"
)

func main() {

	cableType := flag.String("cable", "wg", "peer connection cable type")

	flag.Parse()

	b := broker.NewBroker(*cableType)

	b.StartListeners()
}
