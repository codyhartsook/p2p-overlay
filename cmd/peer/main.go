package main

import (
	"flag"
	"p2p-overlay/pkg/peer"
	"runtime"
)

func main() {

	cableType := flag.String("cable", "wg", "peer connection cable type")

	flag.Parse()

	p := peer.NewPeer(*cableType)
	p.RegisterSelf()

	p.SubscribeToOverlayUpdates()

	runtime.Goexit()
}
