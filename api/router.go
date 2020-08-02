package api
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/Riften/libp2p-playground/host"
)
func InitRouter(n *host.Node) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(outPutInfo)
	apiPeer := r.Group("/peer")
	apiPeer.GET("/info", n.ApiPeerInfo)
	apiPeer.POST("/connect", n.ApiPeerConnect)
	apiPeer.GET("/list", n.ApiPeerList)

	apiSpeed := r.Group("/speed")
	apiSpeed.POST("/send", n.ApiSpeedSend)
	return r
}

func outPutInfo(c *gin.Context) {
	out := c.Request.Host + ": "+ c.Request.Method + " - " + c.Request.RequestURI + "\n"
	out = out + "\tParams:\n"
	fmt.Println(c.PostForm("id"))

	for _, p := range c.Params {
		out = out + "\t\t" + p.Key + ":\t" + p.Value + "\n"
	}
	fmt.Println(out)
	c.Next()
}