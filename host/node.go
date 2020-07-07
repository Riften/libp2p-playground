package host

import (
	"context"
	"fmt"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/Riften/libp2p-playground/repo"
	"time"
)

const defaultConnTimeout = time.Second *30

type Node struct {
	host p2phost.Host
	cfg *repo.Config
	ctx context.Context
}

func NewNode(ctx context.Context, cfg *repo.Config) (*Node, error) {
	h, err := newHost(ctx, cfg)
	if err != nil {
		fmt.Println("Error when create host: ", err)
		return nil, err
	}
	return &Node{host: h, cfg: cfg, ctx: ctx}, nil
}

func (n *Node) Start() {
	foundPeers := n.initMDNS()
	go func() {
		for{
			select {
			case p:= <-foundPeers:
				if !n.IsConnect(p.ID) {
					err := n.host.Connect(n.ctx, p)
					if err != nil {
						fmt.Println("Error when connect ", p.ID.Pretty(), ": ", err)
					} else {
						fmt.Println("Connect ", p.ID.Pretty())
					}
				}
			}
		}
	}()
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

