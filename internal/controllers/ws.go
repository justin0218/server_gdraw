package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"server_gdraw/internal/services"
	"strconv"
)

type Ws struct {
	Base
}

func (s *Ws) Gdraw(c *gin.Context) {
	room := c.Query("room_id")
	wsconnVal, _ := c.Get("wsconn")
	if wsconn, ok := wsconnVal.(*websocket.Conn); ok {
		uidstr := c.Query("uid")
		token := c.Query("token")
		var (
			uidInt64 int64
		)
		uidInt64, _ = strconv.ParseInt(uidstr, 10, 64)
		wsClient := new(services.Ws)
		wsClient.Conn = wsconn
		wsClient.Room = room
		wsClient.Uid = uidInt64
		wsClient.Token = token
		wsClient.MsgChan = make(chan []byte)
		wsClient.ErrMsgChan = make(chan []byte)
		wsClient.BroadcastChan = make(chan []byte)
		wsClient.BroadcastOutSelfChan = make(chan []byte)
		go wsClient.Read()
		go wsClient.Write()
	}
}
