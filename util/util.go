package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"log"
	"os"
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

func IdentityFromKey(sk crypto.PrivKey, pk crypto.PubKey) (config.Identity, error) {
	ident := config.Identity{}
	skbytes, err := sk.Bytes()
	if err != nil {
		log.Println("Error when transfer sk to bytes: ", err)
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)
	id, err := peer.IDFromPublicKey(pk)
	ident.PeerID = id.Pretty()
	return ident, nil
}

func CheckPermissions(path string) (bool, error) {
	_, err := os.Open(path)
	if os.IsNotExist(err) {
		// repo does not exist yet - don't load plugins, but also don't fail
		return false, nil
	}
	if os.IsPermission(err) {
		// repo is not accessible. error out.
		return false, fmt.Errorf("error opening repository at %s: permission denied", path)
	}

	return true, nil
}

func CreateDirIfNotExist(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	} else {
		return errors.New("not a directory: "+path)
	}
}