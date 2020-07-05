package client

import (
	"github.com/gin-gonic/gin"
	"server_gdraw/internal/models"
	"server_gdraw/internal/services"
	"server_gdraw/pkg/resp"
	"strconv"
	"strings"
)

var GdrawUserController gdrawUserController

func init() {
	GdrawUserController = &gdrawUserContron{}
}

type gdrawUserController interface {
	UserAuth(c *gin.Context)
	UserRegist(c *gin.Context)
	Login(c *gin.Context)
	TokenVerify(c *gin.Context)
}

type gdrawUserContron struct {
}

type GdrawUserAuthReq struct {
	EncryptedData string `json:"encryptedData"`
	Code          string `json:"code"`
	Iv            string `json:"iv"`
}

func (ctr *gdrawUserContron) UserAuth(c *gin.Context) {
	req := GdrawUserAuthReq{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	if req.Code == "" || req.EncryptedData == "" || req.Iv == "" {
		resp.RespParamErr(c)
		return
	}
	uinfo, err := services.GdrawUserer.ParseUserInfo(req.EncryptedData, req.Iv, req.Code)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	resp.RespOk(c, uinfo)
	return
}

type UserRegistReq struct {
	Nickname    string `json:"nickname"`
	VcodeId     string `json:"vcode_id"`
	VcodeAnswer string `json:"vcode_answer"`
}

func (ctr *gdrawUserContron) UserRegist(c *gin.Context) {
	req := UserRegistReq{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	tempCon := req.Nickname
	tempCon = strings.ReplaceAll(tempCon, " ", "")
	tempCon = strings.ReplaceAll(tempCon, "\t", "")
	tempCon = strings.ReplaceAll(tempCon, "\n", "")
	if tempCon == "" {
		resp.RespParamErr(c, "昵称不能为空")
		return
	}
	match := services.GdrawVcodeer.Verify(req.VcodeId, req.VcodeAnswer)
	if !match {
		resp.RespCode(c, resp.RESP_CODE_VCODENOMATCH_ERR, "验证码错误")
		return
	}
	uinfo, err := services.GdrawUserer.RegisterUser(services.GdrawWxUserInfo{
		Nickname: req.Nickname,
	})
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	uinfo.Password, err = services.GdrawUserer.MakePassword(uinfo.Id)
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	userinfoFull := services.GdrawUserInfoFull{}
	userinfoFull.Token, err = services.GdrawUserer.CreateToken(int64(uinfo.Id))
	if err != nil {
		return
	}
	userinfoFull.GdrawUser = uinfo
	resp.RespOk(c, userinfoFull)
	return
}

func (ctr *gdrawUserContron) Login(c *gin.Context) {
	req := models.GdrawUser{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	uinfo, loginSuccess, err := services.GdrawUserer.Login(req.Id, req.Password)
	if err != nil {
		resp.RespCode(c, resp.RESP_CODE_LOGINACCOUNT_ERR, "账号或密码错误")
		return
	}
	if !loginSuccess {
		resp.RespCode(c, resp.RESP_CODE_LOGINACCOUNT_ERR, "账号或密码错误")
		return
	}
	userinfoFull := services.GdrawUserInfoFull{}
	userinfoFull.Token, err = services.GdrawUserer.CreateToken(int64(req.Id))
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	userinfoFull.GdrawUser = uinfo
	resp.RespOk(c, userinfoFull)

	return
}

func (ctr *gdrawUserContron) TokenVerify(c *gin.Context) {
	token := c.Query("token")
	uid := c.Query("uid")
	uidInt64, _ := strconv.ParseInt(uid, 10, 64)
	err := services.GdrawUserer.VerifyTokenFunc(uidInt64, token)
	if err != nil {
		resp.RespCode(c, resp.RESP_CODE_NOAUTH_ERR)
		return
	}
	resp.RespOk(c)
	return
}
