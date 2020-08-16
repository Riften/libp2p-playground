package host

import (
	"context"
	"errors"
	"fmt"
	"github.com/Riften/libp2p-playground/repo"
	"github.com/Riften/libp2p-playground/service"
	"github.com/ipfs/go-ipfs/core/node"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-tcp-transport"
	//"github.com/multiformats/go-multiaddr"
	"github.com/ipfs/go-ipfs/core"
	"log"
	"time"
)

const defaultConnTimeout = time.Second * 10

type Node struct {
	host         p2phost.Host
	cfg          *repo.Config
	speedService *service.SpeedService
	tcpService   *service.TCPService
	ipfs         *core.IpfsNode
	ctx          context.Context
}

func EmptyNode(ctx context.Context) *Node {
	return &Node{
		ctx: ctx,
	}
}

func (n *Node) StartLibp2p(cfg *repo.Config) error {
	if n.host != nil {
		return errors.New("libp2p host already start")
	}

	// 0.0.0.0 will listen on any interface device.
	//sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.Port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	privK, err := crypto.UnmarshalPrivateKey(cfg.PrivKey)
	if err != nil {
		fmt.Println("Error when unmarshal privKey: ", err)
		return err
	}
	var transport libp2p.Option
	switch cfg.Transport {
	case "tcp":
		log.Println("Use TCP Transport")
		transport = libp2p.Transport(tcp.NewTCPTransport)
	case "quic":
		log.Println("Use QUIC Transport")
		transport = libp2p.Transport(libp2pquic.NewTransport)
	default:
		log.Println("Transport not specified. Use TCP by default")
		transport = libp2p.Transport(tcp.NewTCPTransport)
	}
	h, err := libp2p.New(
		n.ctx,
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.Port),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", cfg.Port),
			),
		libp2p.Identity(privK),
		//libp2p.Transport()
		transport,
	)
	if err != nil {
		log.Println("Error when create libp2p node: ", err)
		return err
	}
	fmt.Println("Host start at multiaddress:")
	for _, ma := range h.Addrs() {
		fmt.Println(ma.String())
	}

	// Start services
	speedService := service.NewSpeedService(h, context.Background())
	speedService.Start()
	n.host = h
	n.speedService = speedService
	return nil
}

func (n *Node) StartTcp() error {
	if n.tcpService != nil {
		return errors.New("tcp service already start")
	}

	tcpS := service.NewTCPService(n.ctx)
	n.tcpService = tcpS
	return nil
}

func NewIpfsNode(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		log.Println("Error when open fsrepo: ", err)
		return nil, err
	}
	return core.NewNode(ctx, &node.BuildCfg{
		Online:                      true,
		Repo:                        r,
	})
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

