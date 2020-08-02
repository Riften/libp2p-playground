package util

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func PeerJsonIndent(info *peer.AddrInfo, prefix string, indent string) ([]byte, error) {
	out := make(map[string]interface{})
	out["ID"] = info.ID.Pretty()
	var addrs []string
	for _, a := range info.Addrs {
		addrs = append(addrs, a.String())
	}
	out["Addrs"] = addrs
	return json.MarshalIndent(out, prefix, indent)
}

func BuildPeerInfo(id string, addr []string) (*peer.AddrInfo, error) {
	pid, err := peer.Decode(id)
	if err != nil {
		fmt.Println("Error when decode peer id " + id +": ", err)
		return nil, err
	}
	mAddrs := make([]ma.Multiaddr, 0)
	for _, a := range addr {
		mAddrs = append(mAddrs, ma.StringCast(a))
	}

	return &peer.AddrInfo{
		ID:    pid,
		Addrs: mAddrs,
	}, nil
}
