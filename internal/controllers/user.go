package controllers

import (
	"github.com/gin-gonic/gin"
	"server_gdraw/internal/models/user"
	"server_gdraw/internal/services"
	"server_gdraw/pkg/jwt"
	"server_gdraw/pkg/resp"
	"strconv"
)

type User struct {
	Base
	userService services.User
}

func (s *User) UserAuth(c *gin.Context) {
	req := user.AuthReq{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	if req.Code == "" || req.EncryptedData == "" || req.Iv == "" {
		resp.RespParamErr(c)
		return
	}
	uinfo, err := s.userService.ParseUserInfo(req.EncryptedData, req.Iv, req.Code)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	resp.RespOk(c, uinfo)
	return
}

func (s *User) TokenVerify(c *gin.Context) {
	token := c.Query("token")
	uid := c.Query("uid")
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	uidFromToken, err := jwt.VerifyToken(token)
	if err != nil {
		resp.RespCode(c, resp.RESP_CODE_NOAUTH_ERR)
		return
	}
	if uidFromToken != uidInt64 {
		resp.RespCode(c, resp.RESP_CODE_NOAUTH_ERR)
		return
	}
	resp.RespOk(c)
	return
}
