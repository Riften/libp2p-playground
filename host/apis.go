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
	peerId, ok1 := ctx.GetPostForm("id")
	peerAddr, ok2 := ctx.GetPostForm("addr")
	if !(ok1 && ok2) {
		fmt.Println("Not enough param.")
		ctx.Writer.Write([]byte("Not enough param"))
		ctx.Status(http.StatusBadRequest)
		return
	}

	info, err := util.BuildPeerInfo(peerId, []string{peerAddr})
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

func (n *Node) ApiPeerList(ctx *gin.Context) {
	fmt.Println("API: ApiPeerList")

	peers := n.Peers()
	res := make([]string, 0)
	for _, p := range peers {
		res = append(res, p.Pretty())
	}
	responseJson(ctx, res)
}

func (n *Node) ApiSpeedSend(ctx *gin.Context) {
	fmt.Println("API: ApiSpeedSend")

	peerId, ok := ctx.GetPostForm("peer")
	if !ok {
		fmt.Println("Not enough param.")
		ctx.Writer.Write([]byte("Not enough param"))
		ctx.Status(http.StatusBadRequest)
		return
	}
	err := n.speedService.StartSend(peerId)
	if err != nil {
		fmt.Println("Error when start send task: ", err)
		ctx.Writer.Write([]byte("Error when start send task: " +  err.Error()))
		ctx.Status(http.StatusBadGateway)
		return
	}
	ctx.Writer.Write([]byte("Sending start."))
	ctx.Status(http.StatusOK)
}
