package host

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (n *Node) ApiTcpListen(ctx *gin.Context) {
	portStr, ok := ctx.GetPostForm("port")
	if !ok {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("no port specified"))
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	err = n.tcpService.StartListen(port)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.Writer.Write([]byte("tcp listen start success"))
	ctx.Status(http.StatusOK)
}

func (n *Node) ApiTcpSend(ctx *gin.Context) {
	ip, ok1 := ctx.GetPostForm("ip")
	portStr, ok2 := ctx.GetPostForm("port")
	if !(ok1 && ok2) {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("no ip or port specified"))
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	err = n.tcpService.StartSend(ip, port)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.Writer.Write([]byte("tcp send start success"))
	ctx.Status(http.StatusOK)
}