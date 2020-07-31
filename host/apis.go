package host

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Riften/libp2p-playground/util"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p-core/peer"
	"net/http"
)

func responseJson(ctx *gin.Context, value interface{}) {
	resByte, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		fmt.Println("Error when marshal http response json: ", err)
		ctx.AbortWithStatus(http.StatusBadGateway)
	} else {
		ctx.Writer.Write(resByte)
		ctx.Status(http.StatusOK)
	}
}

func (n *Node) ApiPeerInfo(ctx *gin.Context) {
	fmt.Println("API: ApiPeerInfo")
	info := &peer.AddrInfo{
		ID:    n.host.ID(),
		Addrs: n.host.Addrs(),
	}
	res, err := util.PeerJsonIndent(info, "", "\t")
	if err != nil {
		fmt.Println("Error when marshal peer info: ", err)
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	} else {
		ctx.Writer.Write(res)
		ctx.Status(http.StatusOK)
	}
}

func (n *Node) ApiPeerConnect(ctx *gin.Context) {
	fmt.Println("API: ApiPeerConnect")

	info, err := util.BuildPeerInfo(ctx.Param("id"), []string{ctx.Param("addr")})
	if err != nil {
		ctx.Writer.Write([]byte("Fail to build peer info: " + err.Error()))
		ctx.Status(http.StatusBadGateway)
		return
	}

	connCtx, _ := context.WithTimeout(context.Background(), defaultConnTimeout)
	err = n.host.Connect(connCtx, *info)

	if err != nil {
		ctx.Writer.Write([]byte("Fail to connect with peer: " + err.Error()))
		ctx.Status(http.StatusBadGateway)
		return
	}

	ctx.Writer.Write([]byte("Connect successfully."))
	ctx.Status(http.StatusOK)
}
