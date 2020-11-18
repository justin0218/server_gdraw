package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server_gdraw/pkg/resp"
	"time"
)

func HttpUpgrader() gin.HandlerFunc {
	return func(c *gin.Context) {
		wsu := websocket.Upgrader{
			HandshakeTimeout: time.Duration(time.Second * 30),
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		wsconn, err := wsu.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			resp.RespCode(c, resp.RESP_CODE_INTERNAL_ERR)
			c.Abort()
			return
		}
		c.Set("wsconn", wsconn)
		c.Next()
	}
}
