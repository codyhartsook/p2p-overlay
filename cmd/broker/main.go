package main

import (
	"p2p-overlay/pkg/broker"
)

func main() {
	b := broker.NewBroker()

	b.StartListeners()
}
