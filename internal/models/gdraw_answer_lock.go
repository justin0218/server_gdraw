package models

import (
	"github.com/jinzhu/gorm"
	"server_gdraw/api"
	"time"
)

var GdrawAnswerLocker gdrawAnswerLocker

func init() {
	GdrawAnswerLocker = &gdrawAnswerLock{
		TableName: "gdraw_answer_locks",
	}
	GdrawAnswerLocker.InitAllLock()
	//GdrawAnswerLocker.InitLock("1")
	//fmt.Println(GdrawAnswerLocker.Lock("1"))
	//fmt.Println(GdrawAnswerLocker.UnLock("1"))
}

type gdrawAnswerLock struct {
	TableName string
}

type gdrawAnswerLocker interface {
	Lock(roomId string) (b bool)
	UnLock(roomId string) (b bool)
	InitAllLock() (b bool)
	InitLock(roomId string)
}

type GdrawAnswerLock struct {
	RoomId     string    `json:"room_id"`
	IsLock     int       `json:"is_lock"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawAnswerLock) Lock(roomId string) (b bool) {
	err := api.Mysql.Get().Table(this.TableName).Where("room_id = ?", roomId).UpdateColumn("is_lock", gorm.Expr("is_lock - ?", 1)).Error
	if err == nil {
		return true
	}
	return
}

func (this *gdrawAnswerLock) UnLock(roomId string) (b bool) {
	err := api.Mysql.Get().Table(this.TableName).Where("room_id = ?", roomId).Update("is_lock", 1).Error
	if err == nil {
		return true
	}
	return
}

func (this *gdrawAnswerLock) InitLock(roomId string) {
	ret := GdrawAnswerLock{
		RoomId: roomId,
		IsLock: 1,
	}
	api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&ret)
	return
}

func (this *gdrawAnswerLock) InitAllLock() (b bool) {
	err := api.Mysql.Get().Table(this.TableName).Update("is_lock", 1).Error
	if err == nil {
		return true
	}
	return
}
