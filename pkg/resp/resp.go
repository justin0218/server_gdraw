package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RESP_CODE_OK           = 200
	RESP_CODE_PARAM_ERR    = 400
	RESP_CODE_NOAUTH_ERR   = 401
	RESP_CODE_INTERNAL_ERR = 500

	RESP_CODE_USERCREATEDROOM_ERR = 201
	RESP_CODE_NOTFOUNDROOM_ERR    = 4040

	RESP_CODE_VCODENOMATCH_ERR = 4041 //验证码不匹配
	RESP_CODE_LOGINACCOUNT_ERR = 4042 //登录账号或密码错误
)

func RespOk(c *gin.Context, args ...interface{}) {
	var data interface{}
	if len(args) != 0 {
		data = args[0]
	}
	c.JSON(http.StatusOK, gin.H{"code": RESP_CODE_OK, "data": data, "msg": "ok"})
	return
}

func RespCode(c *gin.Context, code int, args ...interface{}) {
	var data interface{}
	if len(args) != 0 {
		data = args[0]
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "data": make([]int, 0), "msg": data})
	return
}

func RespParamErr(c *gin.Context, args ...string) {
	msg := "参数错误"
	if len(args) != 0 {
		msg = args[0]
	}
	c.JSON(http.StatusOK, gin.H{"code": RESP_CODE_PARAM_ERR, "msg": msg})
	return
}

func RespInternalErr(c *gin.Context, args ...string) {
	msg := "内部错误"
	if len(args) != 0 {
		msg = args[0]
	}
	c.JSON(http.StatusOK, gin.H{"code": RESP_CODE_INTERNAL_ERR, "msg": msg})
	return
}
