package host

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

//interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

//Initialize the MDNS service
func (n *Node)initMDNS() chan peer.AddrInfo {
	// An hour might be a long long period in practical applications. But this is fine for us
	ser, err := discovery.NewMdnsService(n.ctx, n.host, 10* time.Second, discovery.ServiceTag)
	if err != nil {
		fmt.Println("Error when init mdns service")
		panic(err)
	}

	//register with service so that we get notified about peer discovery
	noti := &discoveryNotifee{}
	noti.PeerChan = make(chan peer.AddrInfo)

	ser.RegisterNotifee(noti)
	return noti.PeerChan
}