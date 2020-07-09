package host

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/Riften/libp2p-playground/repo"
	"github.com/multiformats/go-multiaddr"
	"time"
)

const defaultConnTimeout = time.Second *30

type Node struct {
	host p2phost.Host
	cfg *repo.Config
	ctx context.Context
}

func NewNode(ctx context.Context, cfg *repo.Config) (*Node, error) {
	h := newHost(cfg.Port)
	fmt.Printf("Host start at multiaddress: /ip4/0.0.0.0/tcp/%d/p2p/%s\n", cfg.Port, h.ID().Pretty())
	return &Node{host: h, cfg: cfg, ctx: ctx}, nil
}

func Start(port int) {
	//n.host.SetStreamHandler()
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

	foundPeers := initMDNS(ctx, host, rendezvous)
	go func() {
		for{
			p:= <-foundPeers
			err := host.Connect(ctx, p)
			if err != nil {
				fmt.Println("Error when connect ", p.ID.Pretty(), ": ", err)
			} else {
				fmt.Println("Connect ", p.ID.Pretty())
			}
		}
	}()

	<-context.Background().Done()
}

func (n *Node) IsConnect(pid peer.ID) bool {
	plist := n.Peers()
	for _, p := range plist {
		if p == pid {
			return true
		}
	}
	return false
}

func (n *Node) Peers() ([]peer.ID) {
	conns := n.host.Network().Conns()
	var result []peer.ID
	for _, c := range conns {
		pid := c.RemotePeer()
		result = append(result, pid)
	}
	return result
}

