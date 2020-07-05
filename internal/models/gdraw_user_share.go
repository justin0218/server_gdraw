package models

import (
	"server_gdraw/api"
	"time"
)

var GdrawUserShareer gdrawUserShareer

func init() {
	GdrawUserShareer = &gdrawUserShare{
		TableName: "gdraw_user_shares",
	}
}

type gdrawUserShare struct {
	TableName string
}

type gdrawUserShareer interface {
	Create(in GdrawUserShare) (res GdrawUserShare, err error)
}

type GdrawUserShare struct {
	Id         int       `json:"id"`
	Uid        int       `json:"uid"`
	BeUid      int       `json:"be_uid"`
	Power      int       `json:"power"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawUserShare) Create(in GdrawUserShare) (res GdrawUserShare, err error) {
	err = api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&in).Error
	res = in
	return
}
