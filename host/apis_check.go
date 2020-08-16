package host

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (n *Node) ApiCheckLibp2p(ctx *gin.Context) {
	if n.host != nil {
		ctx.Next()
	} else {
		ctx.AbortWithError(http.StatusBadGateway, errors.New("lipp2p host not start"))
	}
}

func (n *Node) ApiCheckTcp(ctx *gin.Context) {
	if n.tcpService != nil {
		ctx.Next()
	} else {
		ctx.AbortWithError(http.StatusBadGateway, errors.New("tcp service not start"))
	}
}
