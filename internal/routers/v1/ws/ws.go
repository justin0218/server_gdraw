package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"server_gdraw/internal/services"
	"strconv"
)

var Ws wser

type ws struct {
}

func init() {
	Ws = &ws{}
}

type wser interface {
	Gdraw(c *gin.Context)
}

func (this *ws) Gdraw(c *gin.Context) {
	room := c.Query("room_id")
	wsconnVal, _ := c.Get("wsconn")
	if wsconn, ok := wsconnVal.(*websocket.Conn); ok {
		uidstr := c.Query("uid")
		token := c.Query("token")
		var (
			uidInt64 int64
		)
		uidInt64, _ = strconv.ParseInt(uidstr, 10, 64)
		wsClient := services.NewGdrawWs(wsconn, room, uidInt64, token)
		go wsClient.Read()
		go wsClient.Write()
	}
}
