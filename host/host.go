package host

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
)
/*
func newHost(ctx context.Context,  cfg *repo.Config) (p2phost.Host, error){
	opts, err := option(cfg)
	if err != nil {
		fmt.Printf("Error occurs when build libp2p options\n%s\n", err.Error())
		return nil, err
	}
	return libp2p.New(ctx, opts...)
}
 */

func newHost(port int) p2phost.Host{
	ctx := context.Background()
	r := rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}
	return host
}
