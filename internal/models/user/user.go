package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Model struct {
	Db   *gorm.DB
	Name string
}

func NewModel(db *gorm.DB) *Model {
	return &Model{
		Db:   db,
		Name: "users",
	}
}

type User struct {
	Id        int       `json:"id"`
	Password  string    `json:"password"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Openid    string    `json:"openid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserFull struct {
	Token string `json:"token"`
	User
}

type AuthReq struct {
	EncryptedData string `json:"encryptedData"`
	Code          string `json:"code"`
	Iv            string `json:"iv"`
}

type RegistReq struct {
	Nickname    string `json:"nickname"`
	VcodeId     string `json:"vcode_id"`
	VcodeAnswer string `json:"vcode_answer"`
}

func (s *Model) Create(in User) (out User, err error) {
	db := s.Db.Table(s.Name)
	err = db.Create(&in).Error
	out = in
	return
}

func (s *Model) GetWithOpenid(openid string) (out User, err error) {
	db := s.Db.Table(s.Name)
	err = db.Where("openid = ?", openid).First(&out).Error
	return
}

func (s *Model) GetWithIds(uids []int64) (out []User, err error) {
	if len(uids) == 0 {
		return
	}
	db := s.Db.Table(s.Name)
	err = db.Where("id in (?)", uids).Find(&out).Error
	return
}

func (s *Model) GetWithId(uid int) (out User, err error) {
	db := s.Db.Table(s.Name)
	err = db.Where("id = ?", uid).First(&out).Error
	return
}
