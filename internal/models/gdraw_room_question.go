package models

import (
	"server_gdraw/api"
	"time"
)

var GdrawRoomQuestioner gdrawRoomQuestioner

func init() {
	GdrawRoomQuestioner = &gdrawRoomQuestion{
		TableName: "gdraw_room_questions",
	}
}

type gdrawRoomQuestion struct {
	TableName string
}

type gdrawRoomQuestioner interface {
	Create(in GdrawRoomQuestion) (ret GdrawRoomQuestion, err error)
	GetWithRoomStatus(roomId string, status int) (ret GdrawRoomQuestion, err error)
	UpdateWithId(id int) (err error)
}

type GdrawRoomQuestion struct {
	Id           int       `json:"id"`
	RoomId       string    `json:"room_id"`
	Question     string    `json:"question"`
	QuestionTips string    `json:"question_tips"`
	QuestionId   int       `json:"question_id"`
	Status       int       `json:"status"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

func (this *gdrawRoomQuestion) Create(in GdrawRoomQuestion) (ret GdrawRoomQuestion, err error) {
	err = api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&in).Error
	ret = in
	return
}

func (this *gdrawRoomQuestion) GetWithRoomStatus(roomId string, status int) (ret GdrawRoomQuestion, err error) {
	api.Mysql.Get().Table(this.TableName).Where("room_id = ?", roomId).Where("status = ?", status).First(&ret)
	return
}

func (this *gdrawRoomQuestion) UpdateWithId(id int) (err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("id = ?", id).Update("status", 1).Error
	return
}
