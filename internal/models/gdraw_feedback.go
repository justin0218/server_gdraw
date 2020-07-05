package models

import (
	"server_gdraw/api"
	"time"
)

var GdrawFeedbacker gdrawFeedbacker

func init() {
	GdrawFeedbacker = &gdrawFeedback{
		TableName: "gdraw_feedbacks",
	}
}

type gdrawFeedback struct {
	TableName string
}

type gdrawFeedbacker interface {
	Create(in GdrawFeedback) (ret GdrawFeedback, err error)
}

type GdrawFeedback struct {
	Id         int       `json:"id"`
	Content    string    `json:"content"`
	Uid        int       `json:"uid"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawFeedback) Create(in GdrawFeedback) (ret GdrawFeedback, err error) {
	err = api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&in).Error
	ret = in
	return
}
