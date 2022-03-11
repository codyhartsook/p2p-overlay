package fake

import (
	"context"
	"net"

	pb "p2p-overlay/pkg/grpc"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Fake struct{}

func NewFake() (*Fake, error) {
	return &Fake{}, nil
}

func (w *Fake) Init() error {
	return nil
}

func (w *Fake) AddrAdd() {

}

func (w *Fake) GetPeerTopology() (net.IP, []net.IP, error) {
	return nil, nil, nil
}

func (w *Fake) GetPubKey() string {
	return ""
}

func (w *Fake) RegisterPeer(ctx context.Context, peer wgtypes.PeerConfig) error {
	return nil
}

func (w *Fake) SetAddress(addr string) {
}

func (w *Fake) GetLocalConfig() wgtypes.PeerConfig {
	return wgtypes.PeerConfig{}
}

func (w *Fake) GetPeers(ctx context.Context) ([]wgtypes.PeerConfig, error) {
	return nil, nil
}

func (w *Fake) SyncPeers(ctx context.Context, peers []wgtypes.PeerConfig) error {
	return nil
}

func (w *Fake) ProtobufToPeerConfig(peer *pb.Peer) (wgtypes.PeerConfig, error) {
	return wgtypes.PeerConfig{}, nil
}

func (w *Fake) PeerConfigToProtobuf(conf wgtypes.PeerConfig) (*pb.Peer, error) {
	return &pb.Peer{}, nil
}

func (w *Fake) DeletePeer(ctx context.Context, publicKey string) error {
	return nil
}
