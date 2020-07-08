package api

import (
	"github.com/Riften/libp2p-playground/host"
)

const localhost = "http://localhost"

type Api struct {
	Node *host.Node
	Port int
}

