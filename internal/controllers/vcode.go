package controllers

import (
	"github.com/gin-gonic/gin"
	"server_gdraw/pkg/resp"
	"server_gdraw/pkg/vcode"
)

type Vcode struct {
	Base
}

func (s *Vcode) Get(c *gin.Context) {
	id, b64s, err := vcode.Get()
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
