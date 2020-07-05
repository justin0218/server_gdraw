package client

import (
	"github.com/gin-gonic/gin"
	"server_gdraw/internal/services"
	"server_gdraw/pkg/resp"
)

var GdrawVcodeController gdrawVcodeController

func init() {
	GdrawVcodeController = &gdrawVcodeContron{}
}

type gdrawVcodeController interface {
	Get(c *gin.Context)
	Verify(c *gin.Context)
}

type gdrawVcodeContron struct {
}

func (ctr *gdrawVcodeContron) Get(c *gin.Context) {
	id, b64s, err := services.GdrawVcodeer.Get()
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	res := make(map[string]string)
	res["id"] = id
	res["b64s"] = b64s
	resp.RespOk(c, res)
	return
}

func (ctr *gdrawVcodeContron) Verify(c *gin.Context) {
	id, b64s, err := services.GdrawVcodeer.Get()
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	res := make(map[string]string)
	res["id"] = id
	res["b64s"] = b64s
	resp.RespOk(c, res)
	return
}
