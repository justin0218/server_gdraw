package models

import (
	"fmt"
	"math/rand"
	"server_gdraw/api"
	"server_gdraw/pkg/tools"
	"time"
)

var GdrawUserer gdrawUserer

func init() {
	GdrawUserer = &gdrawUser{
		TableName: "gdraw_users",
	}
}

type gdrawUser struct {
	TableName string
}

type gdrawUserer interface {
	CreateGdrawUser(user GdrawUser) (userinfo GdrawUser, err error)
	FindGdrawUserWithOpenid(openid string) (userinfo GdrawUser, err error)
	FindGdrawUserWithIds(uids []int64) (userinfo []GdrawUser, err error)
	FindGdrawUserWithId(uid int) (userinfo GdrawUser, err error)
	MakePassword(uid int) (password string, err error)
	FindUserById(uid int) (userinfo GdrawUser, err error)
}

type GdrawUser struct {
	Id         int       `json:"id"`
	Password   string    `json:"password"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	Openid     string    `json:"openid"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawUser) CreateGdrawUser(user GdrawUser) (userinfo GdrawUser, err error) {
	err = api.Mysql.Get().Table(this.TableName).Omit("create_time", "update_time").Create(&user).Error
	userinfo = user
	return
}

func (this *gdrawUser) FindGdrawUserWithOpenid(openid string) (userinfo GdrawUser, err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("openid = ?", openid).First(&userinfo).Error
	return
}

func (this *gdrawUser) FindGdrawUserWithIds(uids []int64) (userinfo []GdrawUser, err error) {
	if len(uids) == 0 {
		return
	}
	err = api.Mysql.Get().Table(this.TableName).Where("id in (?)", uids).Find(&userinfo).Error
	return
}

func (this *gdrawUser) FindGdrawUserWithId(uid int) (userinfo GdrawUser, err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("id = ?", uid).First(&userinfo).Error
	return
}

func (this *gdrawUser) MakePassword(uid int) (password string, err error) {
	password = fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	err = api.Mysql.Get().Table(this.TableName).Where("id = ?", uid).Update(&GdrawUser{
		Password: tools.FastFunc.MakeMd5(password),
	}).Error
	return
}

func (this *gdrawUser) FindUserById(uid int) (userinfo GdrawUser, err error) {
	err = api.Mysql.Get().Table(this.TableName).Where("id = ?", uid).First(&userinfo).Error
	return
}
