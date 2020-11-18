package question

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
		Name: "questions",
	}
}

type Question struct {
	Id           int       `json:"id"`
	Question     string    `json:"question"`
	QuestionTips string    `json:"question_tips"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s *Model) GetAll() (ret []Question, err error) {
	db := s.Db.Table(s.Name)
	err = db.Find(&ret).Error
	return
}
