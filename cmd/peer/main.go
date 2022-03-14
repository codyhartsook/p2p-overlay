package main

import (
	"flag"
	"p2p-overlay/pkg/peer"
	"runtime"
)

func main() {

	cableType := flag.String("cable", "wireguard", "peer connection cable type")
	brokerHost := flag.String("broker", "localhost", "broker host")
	zone := flag.String("zone", "", "geographical zone of host")

	flag.Parse()

	p := peer.NewPeer(*cableType, *brokerHost, *zone)
	p.RegisterSelf()

	runtime.Goexit()
}
