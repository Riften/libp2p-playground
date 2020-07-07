package host

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	secio "github.com/libp2p/go-libp2p-secio"
	tls "github.com/libp2p/go-libp2p-tls"
	"github.com/Riften/libp2p-playground/repo"
	"github.com/multiformats/go-multiaddr"
)

// option() build the libp2p options
// options include:
//		- private key, used to set the host identity
//		- protector, used to run host in private network
func option(cfg *repo.Config) ([]libp2p.Option, error){
	var privKey crypto.PrivKey
	var err error

	// load privKey
	opts := make([]libp2p.Option,0)
	if cfg.PrivKey!= nil {
		privKey, err = crypto.UnmarshalPrivateKey(cfg.PrivKey)
		if err != nil {
			fmt.Printf("Error occurs when ummarshal private key from config.\n%s\nCreate host with random key.\n", err.Error())
		} else {
			opts = append(opts, libp2p.Identity(privKey))
		}
	}

	// set security protocol
	// opts = append(opts, libp2p.ChainOptions(libp2p.Security(secio.ID, secio.New), libp2p.Security(tls.ID, tls.New)))
	sourceMultiAddr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/40102")
	opts = append(opts, libp2p.ListenAddrs(sourceMultiAddr))
	return opts, nil
}
