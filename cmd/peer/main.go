package main

import (
	"fmt"
	"p2p-overlay/pkg/peer"
)

func main() {
	p := peer.NewPeer()

	p.RegisterSelf()
	fmt.Println("peer registered...")
}
