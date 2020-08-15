package service

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
)

const GossipTestProtocol = "speedtest/1.0"

type GossipService struct {
	host host.Host
	printer *speedPrinter
	ctx context.Context
}

func NewGossipService(h host.Host, ctx context.Context,) *GossipService {
	newService := &GossipService{
		host:    h,
		printer: &speedPrinter{
			sendSpeedChan: make(chan *record, 5),
			recvSpeedChan: make(chan *record, 5),
		},
		ctx:     ctx,
	}

	return newService
}
