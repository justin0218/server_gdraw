package user_asset

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
		Name: "user_assets",
	}
}

type UserAsset struct {
	Uid       int       `json:"uid"`
	Power     int       `json:"power"`
	Lx        int       `json:"lx"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Model) GetOrCreateWithUid(uid int) (res UserAsset, err error) {
	err = s.Db.Table(s.Name).Where("uid = ?", uid).First(&res).Error
	if err == gorm.ErrRecordNotFound {
		res.Uid = uid
		res.Power = 3600
		res.Lx = 0
		err = s.Db.Table(s.Name).Create(&res).Error
		return
	}
	return
}

func (s *Model) SubPowerWithUid(uid, power int) (err error) {
	db := s.Db.Table(s.Name)
	err = db.Where("uid = ?", uid).UpdateColumn("power", gorm.Expr("power - ?", power)).Error
	return
}

func (s *Model) AddPowerWithUid(uid, power int) (err error) {
	db := s.Db.Table(s.Name)
	err = db.Where("uid = ?", uid).UpdateColumn("power", gorm.Expr("power + ?", power)).Error
	return
}

func (s *Model) UpdateLxWithUid(uid, lx int) (err error) {
	db := s.Db.Table(s.Name)
	err = db.Where("uid = ?", uid).UpdateColumn("lx", gorm.Expr("lx + ?", lx)).Error
	return
}

func (s *Model) GetWithLimit(limit int) (res []UserAsset, err error) {
	db := s.Db.Table(s.Name)
	err = db.Order("lx desc").Limit(limit).Find(&res).Error
	return
}
