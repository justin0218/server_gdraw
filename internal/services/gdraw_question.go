package services

import (
	"encoding/json"
	"math/rand"
	"server_gdraw/api"
	"server_gdraw/internal/models"
	"time"
)

var GdrawQuestioner gdrawQuestioner

func init() {
	GdrawQuestioner = &gdrawQuestion{}
	rand.Seed(time.Now().UnixNano())
}

type gdrawQuestion struct {
	TableName string
}

type gdrawQuestioner interface {
	GetRandQuestions() (ret []models.GdrawQuestion, err error)
	GetQuestionIds() (ret []models.GdrawQuestion, err error)
}

type GdrawQuestion struct {
	Id         int       `json:"id"`
	Question   string    `json:"question"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (this *gdrawQuestion) GetQuestionIds() (ret []models.GdrawQuestion, err error) {
	rk := api.Rds.GetQuestionIdsKey()
	idsStrData, _ := api.Rds.Get().Get(rk).Result()
	if idsStrData != "" {
		err = json.Unmarshal([]byte(idsStrData), &ret)
		if err == nil {
			return
		}
	}
	ret, err = models.GdrawQuestioner.GetQuesions()
	if err != nil {
		return
	}
	questionsBytes, _ := json.Marshal(ret)
	api.Rds.Get().Set(rk, questionsBytes, -1)
	return
}

func (this *gdrawQuestion) GetRandQuestions() (ret []models.GdrawQuestion, err error) {
	questions, _ := this.GetQuestionIds()
	randMap := make(map[int]int)
	for i := 0; i < 100; i++ {
		dx := rand.Intn(len(questions))
		if _, ok := randMap[questions[dx].Id]; !ok {
			ret = append(ret, questions[dx])
			randMap[questions[dx].Id] = 1
			if len(ret) == QUESTIONNUM {
				break
			}
		}
	}
	return
}
