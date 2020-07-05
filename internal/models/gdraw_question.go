package models

import (
	"server_gdraw/api"
	"time"
)

var GdrawQuestioner gdrawQuestioner

func init() {
	GdrawQuestioner = &gdrawQuestion{
		TableName: "gdraw_questions",
	}
}

type gdrawQuestion struct {
	TableName string
}

type gdrawQuestioner interface {
	GetQuesions() (ret []GdrawQuestion, err error)
}

type GdrawQuestion struct {
	Id           int       `json:"id"`
	Question     string    `json:"question"`
	QuestionTips string    `json:"question_tips"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

func (this *gdrawQuestion) GetQuesions() (ret []GdrawQuestion, err error) {
	err = api.Mysql.Get().Table(this.TableName).Find(&ret).Error
	return
}
