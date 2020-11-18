package controllers

import (
	"github.com/gin-gonic/gin"
	"server_gdraw/internal/models/user_share"
	"server_gdraw/internal/services"
	"server_gdraw/pkg/resp"
	"strconv"
)

type Room struct {
	Base
	roomService      services.Room
	indexService     services.Index
	rankService      services.Rank
	userShareService services.UserShare
}

func (s *Room) GetRoomList(c *gin.Context) {
	roomInfo := s.roomService.GetRoomList()
	resp.RespOk(c, roomInfo)
	return
}

func (s *Room) GetIndexData(c *gin.Context) {
	uidStr := c.Query("uid")
	uidInt, _ := strconv.Atoi(uidStr)
	res, err := s.indexService.GetIndexData(uidInt)
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	resp.RespOk(c, res)
	return
}

func (s *Room) GetRankData(c *gin.Context) {
	res, err := s.rankService.GetRankData()
	if err != nil {
		resp.RespInternalErr(c, err.Error())
		return
	}
	resp.RespOk(c, res)
	return
}

func (s *Room) CreateShare(c *gin.Context) {
	req := user_share.UserShare{}
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
	_, err = s.userShareService.Create(req)
	if err != nil {
		resp.RespOk(c)
		return
	}
	resp.RespOk(c)
	return
}
