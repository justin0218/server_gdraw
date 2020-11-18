package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server_gdraw/pkg/jwt"
	"server_gdraw/pkg/resp"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//if new(store.Config).Get().Runmode == "debug" {
		//	c.Next()
		//	return
		//}
		token := c.GetHeader("Authorization")
		uid, err := jwt.VerifyToken(token)
		if err != nil {
			resp.RespCode(c, http.StatusUnauthorized, "未授权")
			c.Abort()
			return
		}
		c.Set("uid", uid)
		c.Next()
	}
}
