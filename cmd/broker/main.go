package main

import (
	"flag"
	"p2p-overlay/pkg/broker"
)

func main() {

	cableType := flag.String("cable", "wireguard", "peer connection cable type")
	natsHost := flag.String("nats", "localhost", "nats host")
	zone := flag.String("zone", "", "geographical zone of host")

	flag.Parse()

	b := broker.NewBroker(*cableType, *natsHost, *zone)

	b.StartListeners()
}
