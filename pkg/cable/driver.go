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
	GetLocalConfig() wgtypes.PeerConfig

	RegisterPeer(ctx context.Context, peer wgtypes.PeerConfig) error

	SetAddress(addr string)

	GetPeers(ctx context.Context) ([]wgtypes.PeerConfig, error)

	SyncPeers(ctx context.Context, peers []wgtypes.PeerConfig) error

	ProtobufToPeerConfig(peer *pb.Peer) (wgtypes.PeerConfig, error)

	PeerConfigToProtobuf(conf wgtypes.PeerConfig) (*pb.Peer, error)

	DeletePeer(ctx context.Context, publicKey string) error

	GetPubKey() string

	Init() error

	AddrAdd()
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
