package cable

import (
	"context"

	pb "p2p-overlay/pkg/grpc"

	log "github.com/sirupsen/logrus"

	F "p2p-overlay/pkg/cable/fake"
	W "p2p-overlay/pkg/cable/wireguard"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Cable interface {
	RegisterPeer(ctx context.Context, peer wgtypes.PeerConfig) error
	GetPeers(ctx context.Context) ([]wgtypes.PeerConfig, error)
	SyncPeers(ctx context.Context, peers []wgtypes.PeerConfig) error
	DeletePeer(ctx context.Context, publicKey string) error
	Init() error
}

const (
	wg   = "wireguard"
	fake = "fake"
)

func NewCable(cableType string) Cable {
	switch cableType {
	case wg:
		w, err := W.NewWGCtrl()
		if err != nil {
			log.Fatalf("error creating wireguard cable: %v", err)
		}
		return w
	case fake:
		f, _ := F.NewFake()
		return f
	default:
		log.Fatalf("driver type not matched: %s", cableType)
	}

	return nil
}

func ProtobufPeerToConfig(peer *pb.Peer) wgtypes.PeerConfig {
	return wgtypes.PeerConfig{}
}
