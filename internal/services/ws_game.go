package services

import (
	"context"
	"server_gdraw/internal/models/question"
	"server_gdraw/internal/models/room_question"
	"time"
)

type GAME_READY_START_DATA struct {
	Event int `json:"event"`
	Data  struct {
		CountdownSecond int    `json:"countdown_second"`
		Question        string `json:"question"`
		QuestionTips    string `json:"question_tips"`
		Status          int    `json:"status"`
	} `json:"data"`
}

type EVENT_GAME_ANSWER_DATA struct {
	Event int `json:"event"`
	Data  struct {
		Avatar       string `json:"avatar"`
		Nickname     string `json:"nickname"`
		Answer       string `json:"answer"`
		AnswerResult int    `json:"answer_result"` //1.作答正确，0不正确
		Sec          int    `json:"sec"`
		Uid          int    `json:"uid"`
		Room         string `json:"room"`
	} `json:"data"`
}

func (s *Ws) countdowner(ctx context.Context, countdownSec int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case c <- countdownSec:
				countdownSec--
			}
			time.Sleep(time.Second * 1)
		}
	}()
	return c
}

func (s *Ws) questionEmiter(roomId string) (question question.Question) {
	roomQuestion, _ := room_question.NewModel(s.mysql.Get()).GetWithRoomId(roomId, 0)
	question.Id = roomQuestion.QuestionId
	question.Question = roomQuestion.Question
	question.QuestionTips = roomQuestion.QuestionTips
	_ = room_question.NewModel(s.mysql.Get()).UpdateWithId(roomQuestion.Id)
	return
}
