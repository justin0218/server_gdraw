package user_share

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
		Name: "user_shares",
	}
}

type UserShare struct {
	Id        int       `json:"id"`
	Uid       int       `json:"uid"`
	BeUid     int       `json:"be_uid"`
	Power     int       `json:"power"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Model) Create(in UserShare) (res UserShare, err error) {
	db := s.Db.Table(s.Name)
	err = db.Create(&in).Error
	res = in
	return
}
