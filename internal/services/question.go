package services

import (
	"encoding/json"
	"math/rand"
	"server_gdraw/internal/models/question"
)

type Question struct {
	base
}

func (s *Question) GetQuestionIds() (ret []question.Question, err error) {
	rk := s.redis.GetQuestionIdsKey()
	idsStrData, _ := s.redis.Get().Get(rk).Result()
	if idsStrData != "" {
		err = json.Unmarshal([]byte(idsStrData), &ret)
		if err == nil {
			return
		}
	}
	ret, err = question.NewModel(s.mysql.Get()).GetAll()
	if err != nil {
		return
	}
	questionsBytes, _ := json.Marshal(ret)
	s.redis.Get().Set(rk, questionsBytes, -1)
	return
}

func (s *Question) GetRandQuestions() (ret []question.Question, err error) {
	questions, _ := s.GetQuestionIds()
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
