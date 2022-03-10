package main

import (
	"flag"
	"p2p-overlay/pkg/peer"
	"runtime"
)

func main() {

	cableType := flag.String("cable", "wireguard", "peer connection cable type")
	brokerHost := flag.String("broker", "localhost", "broker host")

	flag.Parse()

	p := peer.NewPeer(*cableType, *brokerHost)
	p.RegisterSelf()

	p.SubscribeToOverlayUpdates()

	runtime.Goexit()
}
