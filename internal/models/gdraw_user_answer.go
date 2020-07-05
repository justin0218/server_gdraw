package models

import (
	"server_gdraw/api"
	"time"
)

var GdrawUserAnswerer gdrawUserAnswerer

func init() {
	GdrawUserAnswerer = &gdrawUserAnswer{
		TableName: "gdraw_user_answers",
	}
}

type gdrawUserAnswer struct {
	TableName string
}

type gdrawUserAnswerer interface {
	Create(in GdrawUserAnswer) (err error)
}

type GdrawUserAnswer struct {
	Id         int       `json:"id"`
	Uid        int       `json:"uid"`
	AnswerJson string    `json:"answer_json"`
	IsRight    int       `json:"is_right"`
	QuestionId int       `json:"question_id"`
	RoomId     string    `json:"room_id"`
	AnswerTime int       `json:"answer_time"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawUserAnswer) Create(in GdrawUserAnswer) (err error) {
	err = api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&in).Error
	return
}
