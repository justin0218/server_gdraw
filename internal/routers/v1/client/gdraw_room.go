package client

import (
	"github.com/gin-gonic/gin"
	"server_gdraw/internal/models"
	"server_gdraw/internal/services"
	"server_gdraw/pkg/resp"
	"strconv"
)

var GdrawController gdrawController

func init() {
	GdrawController = &gdrawContron{}
}

type gdrawController interface {
	GetRoomList(c *gin.Context)
	GetIndexData(c *gin.Context)
	GetRankData(c *gin.Context)
	CreateFeedback(c *gin.Context)
	CreateShare(c *gin.Context)
}

type gdrawContron struct {
}

func (ctr *gdrawContron) GetRoomList(c *gin.Context) {
	roomInfo := services.GdrawRoomer.GetRoomList()
	resp.RespOk(c, roomInfo)
	return
}

func (ctr *gdrawContron) GetIndexData(c *gin.Context) {
	uidStr := c.Query("uid")
	uidInt, _ := strconv.Atoi(uidStr)
	res, err := services.GdrawIndexer.GetIndexData(uidInt)
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	resp.RespOk(c, res)
	return
}

func (ctr *gdrawContron) GetRankData(c *gin.Context) {
	res, err := services.GdrawRanker.GetRankData()
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	resp.RespOk(c, res)
	return
}

func (ctr *gdrawContron) CreateFeedback(c *gin.Context) {
	req := models.GdrawFeedback{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.RespParamErr(c, err.Error())
		return
	}
	if req.Content == "" || req.Uid <= 0 {
		resp.RespParamErr(c, "参数错误")
		return
	}
	_, err = services.GdrawFeedbacker.Create(req)
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	resp.RespOk(c)
	return
}

func (ctr *gdrawContron) CreateShare(c *gin.Context) {
	req := models.GdrawUserShare{}
	err := c.BindJSON(&req)
	if err != nil {
		resp.RespOk(c)
		return
	}
	if req.Uid == 0 {
		resp.RespOk(c)
		return
	}
	if req.Uid == req.BeUid {
		resp.RespOk(c)
		return
	}
	_, err = services.GdrawUserShareer.Create(req)
	if err != nil {
		resp.RespOk(c)
		return
	}
	resp.RespOk(c)
	return
}
