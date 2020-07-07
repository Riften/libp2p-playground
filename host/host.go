package host

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	"github.com/Riften/libp2p-playground/repo"
)

func newHost(ctx context.Context,  cfg *repo.Config) (p2phost.Host, error){
	opts, err := option(cfg)
	if err != nil {
		fmt.Printf("Error occurs when build libp2p options\n%s\n", err.Error())
		return nil, err
	}
	return libp2p.New(ctx, opts...)
}

