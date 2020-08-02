package host

import (
	"context"
	"fmt"
	"github.com/Riften/libp2p-playground/repo"
	"github.com/Riften/libp2p-playground/service"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"time"
)

const defaultConnTimeout = time.Second * 10

type Node struct {
	host p2phost.Host
	cfg *repo.Config
	speed *service.SpeedService
	ctx context.Context
}

func NewNode(ctx context.Context, cfg *repo.Config) (*Node, error) {
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.Port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	privK, err := crypto.UnmarshalPrivateKey(cfg.PrivKey)
	if err != nil {
		fmt.Println("Error when unmarshal privKey: ", err)
		return nil, err
	}
	h, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privK),
	)
	fmt.Printf("Host start at multiaddress: /ip4/0.0.0.0/tcp/%d/p2p/%s\n", cfg.Port, h.ID().Pretty())

	// Start services
	speedService := service.NewSpeedService(h, context.Background())
	speedService.Start()

	return &Node{host: h, cfg: cfg, speed: speedService, ctx: ctx}, nil
}

func (n *Node) Host() (p2phost.Host) {
	return n.host
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

