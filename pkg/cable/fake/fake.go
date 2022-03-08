package fake

import (
	"context"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Fake struct{}

func NewFake() (*Fake, error) {
	return &Fake{}, nil
}

func (w *Fake) Init() error {
	return nil
}

func (w *Fake) RegisterPeer(ctx context.Context, deviceName string, peer wgtypes.PeerConfig) error {
	return nil
}

func (w *Fake) GetPeers(ctx context.Context) ([]wgtypes.Peer, error) {
	return nil, nil
}

func (w *Fake) SyncPeers(ctx context.Context, peers []wgtypes.Peer) error {
	return nil
}

func (w *Fake) DeletePeer(ctx context.Context, deviceName string, publicKey string) error {
	return nil
}
