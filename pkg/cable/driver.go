package cable

import (
	"context"
	"log"

	W "p2p-overlay/pkg/cable/wireguard"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Cable interface {
	RegisterPeers(ctx context.Context, deviceName string, peers []wgtypes.PeerConfig) error
	DeletePeers(ctx context.Context, deviceName string, publicKeys []string) error
	Init() error
}

var cableType = "wireguard"

const (
	wg = "wireguard"
)

func NewCable() Cable {
	switch cableType {
	case wg:
		w, err := W.NewWGCtrl()
		if err != nil {
			log.Fatalf("error creating wireguard cable: %v", err)
		}
		return w
	default:
		log.Fatalf("driver type not matched: %s", wg)
	}

	return nil
}
