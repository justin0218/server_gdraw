package room_question

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
		Name: "room_questions",
	}
}

type RoomQuestion struct {
	Id           int       `json:"id"`
	RoomId       string    `json:"room_id"`
	Question     string    `json:"question"`
	QuestionTips string    `json:"question_tips"`
	QuestionId   int       `json:"question_id"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s *Model) Create(in RoomQuestion) (ret RoomQuestion, err error) {
	db := s.Db.Table(s.Name)
	err = db.Create(&in).Error
	ret = in
	return
}

func (s *Model) GetWithRoomId(roomId string, status int) (ret RoomQuestion, err error) {
	db := s.Db.Table(s.Name)
	db.Where("room_id = ?", roomId).Where("status = ?", status).First(&ret)
	return
}

func (s *Model) UpdateWithId(id int) (err error) {
	db := s.Db.Table(s.Name)
	err = db.Where("id = ?", id).Update("status", 1).Error
	return
}
