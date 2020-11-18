package user_answer

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
		Name: "user_answers",
	}
}

type UserAnswer struct {
	Id         int       `json:"id"`
	Uid        int       `json:"uid"`
	AnswerJson string    `json:"answer_json"`
	IsRight    int       `json:"is_right"`
	QuestionId int       `json:"question_id"`
	RoomId     string    `json:"room_id"`
	AnswerTime int       `json:"answer_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *Model) Create(in UserAnswer) (err error) {
	err = s.Db.Table(s.Name).Create(&in).Error
	return
}
