package models

import (
	"github.com/jinzhu/gorm"
	"server_gdraw/api"
	"time"
)

var GdrawUserAsseter gdrawUserAsseter

func init() {
	GdrawUserAsseter = &gdrawUserAsset{
		TableName: "gdraw_user_assets",
	}
}

type gdrawUserAsset struct {
	TableName string
}

type gdrawUserAsseter interface {
	GetUserAsset(uid int) (res GdrawUserAsset, err error)
	UpdateUserPower(uid, power int) (err error)
	UpdateUserLx(uid, lx int) (err error)
	GetRankData(limit int) (res []GdrawUserAsset, err error)
	AddUserPower(uid, power int) (err error)
}

type GdrawUserAsset struct {
	Uid        int       `json:"uid"`
	Power      int       `json:"power"`
	Lx         int       `json:"lx"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawUserAsset) GetUserAsset(uid int) (res GdrawUserAsset, err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("uid = ?", uid).First(&res).Error
	if err == gorm.ErrRecordNotFound {
		res.Uid = uid
		res.Power = 3600
		res.Lx = 0
		err = api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&res).Error
		return
	}
	return
}

func (this *gdrawUserAsset) UpdateUserPower(uid, power int) (err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("uid = ?", uid).UpdateColumn("power", gorm.Expr("power - ?", power)).Error
	return
}

func (this *gdrawUserAsset) AddUserPower(uid, power int) (err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("uid = ?", uid).UpdateColumn("power", gorm.Expr("power + ?", power)).Error
	return
}

func (this *gdrawUserAsset) UpdateUserLx(uid, lx int) (err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("uid = ?", uid).UpdateColumn("lx", gorm.Expr("lx + ?", lx)).Error
	return
}

func (this *gdrawUserAsset) GetRankData(limit int) (res []GdrawUserAsset, err error) {
	err = api.Mysql.Get().Table(this.TableName).Order("lx desc").Limit(limit).Find(&res).Error
	return
}
