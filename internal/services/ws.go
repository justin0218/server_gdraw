package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"server_gdraw/internal/models"
	"server_gdraw/internal/models/question"
	"server_gdraw/internal/models/room_question"
	"server_gdraw/internal/models/user"
	"server_gdraw/internal/models/user_answer"
	"server_gdraw/internal/models/user_asset"
	"server_gdraw/pkg/jwt"
	"server_gdraw/pkg/tools"
	"server_gdraw/pkg/tools/store"
	"strconv"
	"strings"
	"sync"
)

const (
	EVENT_ERR_OTHER       = -2 //错误，其他错误
	EVENT_MEMBER_LIST     = 0  //成员列表
	EVENT_DRAW_CONTENT    = 1  //画的时候，传的坐标信息
	EVENT_GAME_START      = 2  //游戏开始
	EVENT_GAME_ANSWER     = 3  //作答
	EVENT_RANKING         = 4  //排行榜数据
	EVENT_GAME_STATUS     = 5  //游戏状态，1.游戏中 0未开始
	EVENT_USER_ROLE       = 6  //用户角色  0成员 1.房主
	EVENT_INITPATH_DATA   = 7  //房间的初始动画路径
	EVENT_ROOMLEADER_INFO = 8  //通知房主信息
	EVENT_RESTAGME        = 9  //重新游戏，彻底结束
	EVENT_CLEAR_CANVAS    = 10 //清空画布
	EVENT_ANSWERRIGHT     = 11 //作答正确弹窗
	EVENT_NOTICE_IN       = 12 //通知进入房间
	EVENT_NOTICE_OUT      = 13 //通知离开房间
	EVENT_TIPS            = 14 //tips通知
)

const (
	GAME_START_STATUS_NOCANSTART = 1 //游戏未满足开始条件
	GAME_START_STATUS_CANSTART   = 0 //游戏满足开始条件

	GAME_STARTED       = 1 //游戏进行中
	GAME_WAITING_START = 0 //游戏等待开始
	GAME_END           = 2 //游戏结束
)

const (
	ANSWERTIME  = 60
	QUESTIONNUM = 5
	RANKTIME    = 30
)

type EVENT_CLEAR_CANVAS_DATA struct {
	Event int `json:"event"`
}

type EVENT_NOTICE_INOUT_DATA struct {
	Event int       `json:"event"`
	Uinfo user.User `json:"uinfo"`
}

type EVENT_ANSWERRIGHT_DATA struct {
	Event    int    `json:"event"`
	Question string `json:"question"`
	Sec      int    `json:"sec"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type EVENT_RANKING_DATA struct {
	Event int             `json:"event"`
	Data  models.RankData `json:"data"`
}

type EVENT_RESTAGME_DATA struct {
	Event int `json:"event"`
}

type EVENT_ROOMLEADER_DATA struct {
	Event int       `json:"event"`
	Data  user.User `json:"data"`
}

type EVENT_USER_ROLE_DATA struct { //角色信息
	Event int `json:"event"`
	Role  int `json:"role"`
}

type EVENT_ERR_DATA struct { //错误信息
	Event int    `json:"event"`
	Msg   string `json:"msg"`
}

type EVENT_MEMBER_LIST_DATA struct { //房间成员
	Event int                           `json:"event"`
	Data  []EVENT_MEMBER_LIST_DATA_DATA `json:"data"`
}

type EVENT_GAME_STATUS_DATA struct { //游戏状态
	Event       int `json:"event"`
	RoomStatus  int `json:"room_status"`
	ReadyStatus int `json:"ready_status"`
}

type EVENT_GAME_END_DATA struct { //游戏结束
	Event int `json:"event"`
}

type EVENT_MEMBER_LIST_DATA_DATA struct {
	Role     int       `json:"role"`
	UserInfo user.User `json:"user_info"`
}

type WsBaseMsg struct {
	Event int `json:"event"`
}

type Ws struct {
	Conn                 *websocket.Conn
	Room                 string
	Uid                  int64
	MsgChan              chan []byte
	Token                string
	ErrMsgChan           chan []byte
	BroadcastChan        chan []byte
	BroadcastOutSelfChan chan []byte
	Mutex                sync.Mutex
	base
	roomService Room
}

type Wser interface {
	Read()
	Write()
}

//单独携程写消息
func (s *Ws) Write() {
	defer func() {
		_ = s.Conn.Close()
		//s.Conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(3))
	}()
OUT:
	for {
		select {
		case message, ok := <-s.MsgChan:
			if !ok {
				break OUT
			}
			s.Mutex.Lock()
			_ = s.Conn.WriteMessage(websocket.TextMessage, message)
			s.Mutex.Unlock()
			break
		case message, ok := <-s.ErrMsgChan:
			if !ok {
				break OUT
			}
			if conn, ok := store.UID_WSCONN_MAP.Load(s.Uid); ok {
				conn.(*Ws).Mutex.Lock()
				_ = conn.(*Ws).Conn.WriteMessage(websocket.TextMessage, message)
				conn.(*Ws).Mutex.Unlock()
			}
			break OUT
		case message, ok := <-s.BroadcastChan:
			if !ok {
				break OUT
			}
			members := s.getOnlineRoomMems()
			for _, uidInt64 := range members {
				if conn, ok := store.UID_WSCONN_MAP.Load(uidInt64); ok {
					conn.(*Ws).Mutex.Lock()
					_ = conn.(*Ws).Conn.WriteMessage(websocket.TextMessage, message)
					conn.(*Ws).Mutex.Unlock()
				}
			}
			break
		case message, ok := <-s.BroadcastOutSelfChan:
			if !ok {
				break OUT
			}
			members := s.getOnlineRoomMems()
			for _, uid := range members {
				if uid != s.Uid {
					if conn, ok := store.UID_WSCONN_MAP.Load(uid); ok {
						conn.(*Ws).MsgChan <- message
					}
				}
			}
		}
	}
}

//单独携程读取消息
func (s *Ws) Read() {
	s.inroom()
	defer s.outroom()
	msgchan := make(chan []byte)
	for {
		go func() {
			select {
			case m := <-msgchan:
				wsBaseMsg := WsBaseMsg{}
				err := json.Unmarshal(m, &wsBaseMsg)
				if err != nil {
					s.ErrMsgChan <- []byte(fmt.Sprintf("参数错误:%v", err))
					return
				}
				switch wsBaseMsg.Event {
				case EVENT_DRAW_CONTENT:
					if s.Room != fmt.Sprintf("%d", s.Uid) {
						s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("房主才有权限画画哦"))
						return
					}

					if v, ok := store.ROOM_READY_STATUS.Load(s.Room); ok {
						if v.(EVENT_GAME_STATUS_DATA).RoomStatus != GAME_STARTED {
							s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前游戏状态不能画画哦"))
							return
						}
					} else {
						s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前游戏状态不能画画哦"))
						return
					}

					wsBaseMsg := store.EVENT_DRAW_CONTENT_DATA{}
					_ = json.Unmarshal(m, &wsBaseMsg)
					wsBaseMsg.Event = EVENT_DRAW_CONTENT
					if pathData, ok := store.ROOM_CONTENT_DATA.Load(s.Room); ok {
						pdata := pathData.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
						pdata = append(pdata, wsBaseMsg.Data)
						store.ROOM_CONTENT_DATA.Store(s.Room, pdata)
					} else {
						pdata := []store.EVENT_DRAW_CONTENT_DATA_DATA{}
						pdata = append(pdata, wsBaseMsg.Data)
						store.ROOM_CONTENT_DATA.Store(s.Room, pdata)
					}
					s.BroadcastOutSelfChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
					break
				case EVENT_GAME_START:
					if s.Room != fmt.Sprintf("%d", s.Uid) {
						s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("房主才有权限开始游戏"))
						return
					}
					if v, geted := store.ROOM_READY_STATUS.Load(s.Room); geted {
						if v.(EVENT_GAME_STATUS_DATA).RoomStatus != GAME_WAITING_START && v.(EVENT_GAME_STATUS_DATA).ReadyStatus != GAME_START_STATUS_CANSTART {
							s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前状态不能开始游戏"))
							return
						}
					} else {
						s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前状态不能开始游戏"))
						return
					}

					//models.GdrawAnswerLocker.InitLock(s.Room)

					wsBaseMsg := EVENT_GAME_STATUS_DATA{}
					wsBaseMsg.Event = EVENT_GAME_STATUS
					wsBaseMsg.ReadyStatus = GAME_START_STATUS_CANSTART
					wsBaseMsg.RoomStatus = GAME_STARTED
					store.ROOM_READY_STATUS.Store(s.Room, wsBaseMsg)
					s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)

					//选题入库
					randQs, _ := question.NewModel(s.mysql.Get()).GetAll()
					for _, q := range randQs {
						_, _ = room_question.NewModel(s.mysql.Get()).Create(room_question.RoomQuestion{
							RoomId:       s.Room,
							Question:     q.Question,
							QuestionId:   q.Id,
							QuestionTips: q.QuestionTips,
						})
					}
					s.loopGame()
					break
				case EVENT_GAME_ANSWER:
					if s.Room == fmt.Sprintf("%d", s.Uid) {
						s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("房主不能回答问题哦"))
						return
					}
					wsBaseMsg := EVENT_GAME_ANSWER_DATA{}
					_ = json.Unmarshal(m, &wsBaseMsg)

					tempCon := wsBaseMsg.Data.Answer
					tempCon = strings.ReplaceAll(tempCon, " ", "")
					tempCon = strings.ReplaceAll(tempCon, "\t", "")
					tempCon = strings.ReplaceAll(tempCon, "\n", "")
					if tempCon == "" {
						s.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("请输入合法内容"))
						return
					}
					if v, ok := store.ROOM_ANSWER_MAP.Load(s.Room); ok {
						if v.(question.Question).Question == wsBaseMsg.Data.Answer {
							//locker := models.GdrawAnswerLocker.Lock(s.Room)
							//if !locker {
							//	return
							//}
							//models.GdrawAnswerLocker.UnLock(s.Room)
							store.ROOM_ANSWER_MAP.Delete(s.Room) //删除答案
							if cancel, ok := store.ROOM_TIMER_CANCEL.Load(s.Room); ok {
								cancel.(context.CancelFunc)()
								if glCountDown, oked := store.ROOM_COUNTDOWN.Load(s.Room); oked {
									s.CreateUserAnswerRecord(wsBaseMsg.Data.Uid, v.(question.Question).Id, glCountDown.(int), 1)
								}
								wsBaseMsg.Event = EVENT_ANSWERRIGHT
								wsBaseMsg.Data.AnswerResult = 1
								//展示答对的信息
								cx, cl := context.WithCancel(context.Background())
								showTime := s.countdowner(cx, 4)
								for sec := range showTime {
									if sec < 0 {
										cl()
										break
									}
									wsBaseMsg.Data.Sec = sec
									s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
								}
								s.loopGame()
							}
						} else {
							s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
						}
					} else {
						s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
					}
					break
				case EVENT_CLEAR_CANVAS:
					wsBaseMsg := EVENT_CLEAR_CANVAS_DATA{}
					wsBaseMsg.Event = EVENT_CLEAR_CANVAS
					store.ROOM_CONTENT_DATA.Delete(s.Room)
					s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
					break
				case EVENT_TIPS:
					s.BroadcastChan <- m
					break

				}
			}
		}()
		_, rmsg, err := s.Conn.ReadMessage()
		if err != nil {
			break
		}
		msgchan <- rmsg
	}
}

//通知历史的路径数据
func (s *Ws) noticeRoomPathData() {
	wsBaseMsg := store.EVENT_INITPATH_DATA_DATA{}
	wsBaseMsg.Event = EVENT_INITPATH_DATA
	if v, ok := store.ROOM_CONTENT_DATA.Load(s.Room); ok {
		wsBaseMsg.Data = v.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
	}
	s.MsgChan <- []byte(tools.FastFunc.DataToString(wsBaseMsg))
}

//通知错误
func (s *Ws) noticeSelfErr(eventCode int, errMsg string) {
	wsBaseMsg := EVENT_ERR_DATA{}
	wsBaseMsg.Event = eventCode
	wsBaseMsg.Msg = errMsg
	s.MsgChan <- []byte(tools.FastFunc.DataToString(wsBaseMsg))
}

//通知房间内成员
func (s *Ws) noticeMembersToRoom() {
	wsBaseMsg := EVENT_MEMBER_LIST_DATA{}
	wsBaseMsg.Event = EVENT_MEMBER_LIST
	wsBaseMsg.Data = s.GetMembersWithRoom()
	s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
}

//发送房间状态
func (s *Ws) noticeRoomStatus() {
	if v, ok := store.ROOM_READY_STATUS.Load(s.Room); ok {
		if v.(EVENT_GAME_STATUS_DATA).RoomStatus == GAME_END || v.(EVENT_GAME_STATUS_DATA).RoomStatus == GAME_STARTED {
			s.BroadcastChan <- tools.FastFunc.DataToBytes(v.(EVENT_GAME_STATUS_DATA))
			return
		}
		if v.(EVENT_GAME_STATUS_DATA).RoomStatus == GAME_WAITING_START {
			members := s.getOnlineRoomMems()
			wsBaseMsg := EVENT_GAME_STATUS_DATA{}
			wsBaseMsg.Event = EVENT_GAME_STATUS
			if len(members) > 1 {
				//可以开始游戏
				wsBaseMsg.ReadyStatus = GAME_START_STATUS_CANSTART
			} else {
				wsBaseMsg.ReadyStatus = GAME_START_STATUS_NOCANSTART
			}
			wsBaseMsg.RoomStatus = GAME_WAITING_START
			store.ROOM_READY_STATUS.Store(s.Room, wsBaseMsg)
			s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
		}
		return
	}
	members := s.getOnlineRoomMems()
	wsBaseMsg := EVENT_GAME_STATUS_DATA{}
	wsBaseMsg.Event = EVENT_GAME_STATUS
	if len(members) > 1 {
		//可以开始游戏
		wsBaseMsg.ReadyStatus = GAME_START_STATUS_CANSTART
	} else {
		wsBaseMsg.ReadyStatus = GAME_START_STATUS_NOCANSTART
	}
	wsBaseMsg.RoomStatus = GAME_WAITING_START
	store.ROOM_READY_STATUS.Store(s.Room, wsBaseMsg)
	s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
	return
}

//获取房间成员
func (s *Ws) GetMembersWithRoom() (data []EVENT_MEMBER_LIST_DATA_DATA) {
	uids := s.getOnlineRoomMems()
	membersFull, _ := user.NewModel(s.mysql.Get()).GetWithIds(uids)
	for _, val := range membersFull {
		item := EVENT_MEMBER_LIST_DATA_DATA{}
		item.UserInfo = val
		data = append(data, item)
	}
	return
}

//进入房间
func (s *Ws) inroom() {
	_, err := jwt.VerifyToken(s.Token)
	if err != nil {
		s.MsgChan <- []byte(fmt.Sprintf(`{"event":-1,"msg":"%s"}`, err.Error()))
		return
	}
	//效验体力
	userAsset, _ := user_asset.NewModel(s.mysql.Get()).GetOrCreateWithUid(int(s.Uid))
	if userAsset.Power <= 0 {
		s.MsgChan <- []byte(fmt.Sprintf(`{"event":-1,"msg":"体力值不足"}`))
		return
	}
	err = s.roomService.JoinRoom(s.Uid, s.Room)
	if err != nil {
		s.MsgChan <- []byte(fmt.Sprintf(`{"event":-1,"msg":"%s"}`, err.Error()))
		return
	}
	store.UID_WSCONN_MAP.Store(s.Uid, s)

	s.noticeSelfRole()
	//通知房间内所有成员当前房间状态
	s.noticeRoomStatus()
	//进入房间就需要通知的
	s.noticeRoomLeaderInfo()
	s.noticeMembersToRoom()
	s.noticeRoomPathData()

	noticeMsg := EVENT_NOTICE_INOUT_DATA{}
	noticeMsg.Event = EVENT_NOTICE_IN
	noticeMsg.Uinfo, _ = user.NewModel(s.mysql.Get()).GetWithId(int(s.Uid))
	s.BroadcastChan <- tools.FastFunc.DataToBytes(noticeMsg)
	return
}

func (s *Ws) noticeSelfRole() {
	wsBaseMsg := EVENT_USER_ROLE_DATA{}
	wsBaseMsg.Event = EVENT_USER_ROLE
	if s.Room == fmt.Sprintf("%d", s.Uid) {
		wsBaseMsg.Role = 1
	}
	s.MsgChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
}

//获取房主的信息
func (s *Ws) noticeRoomLeaderInfo() (uninfo user.User) {
	uidInt, _ := strconv.Atoi(s.Room)
	uninfo, _ = user.NewModel(s.mysql.Get()).GetWithId(uidInt)
	wsBaseMsg := EVENT_ROOMLEADER_DATA{}
	wsBaseMsg.Event = EVENT_ROOMLEADER_INFO
	wsBaseMsg.Data = uninfo
	s.MsgChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
	return
}

//离开房间
func (s *Ws) outroom() {
	if _, ok := store.UID_WSCONN_MAP.Load(s.Uid); ok {
		store.UID_WSCONN_MAP.Delete(s.Uid)
		uids := s.getOnlineRoomMems()
		for i, uid := range uids {
			if uid == s.Uid {
				uids = append(uids[:i], uids[i+1:]...)
				break
			}
		}
		store.ROOM_ALL_PERSON.Store(s.Room, uids)
		s.noticeMembersToRoom()
		//通知房间内所有成员当前房间状态
		s.noticeRoomStatus()

		//通知离开房间
		noticeMsg := EVENT_NOTICE_INOUT_DATA{}
		noticeMsg.Event = EVENT_NOTICE_OUT
		noticeMsg.Uinfo, _ = user.NewModel(s.mysql.Get()).GetWithId(int(s.Uid))
		s.BroadcastChan <- tools.FastFunc.DataToBytes(noticeMsg)
	}
}

func (s *Ws) getOnlineRoomMems() (uids []int64) {
	if v, ok := store.ROOM_ALL_PERSON.Load(s.Room); ok {
		uids = v.([]int64)
		return
	}
	return
}

//游戏结束
func (s *Ws) gameEnd() {
	store.ROOM_ANSWER_MAP.Delete(s.Room)
	wsBaseMsg := EVENT_GAME_STATUS_DATA{}
	wsBaseMsg.Event = EVENT_GAME_STATUS
	wsBaseMsg.RoomStatus = GAME_END
	wsBaseMsg.ReadyStatus = GAME_START_STATUS_NOCANSTART
	store.ROOM_READY_STATUS.Store(s.Room, wsBaseMsg)
	s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)

	rankingData := EVENT_RANKING_DATA{}
	rankingData.Event = EVENT_RANKING
	rankingData.Data = s.getRanking()

	lxMap := make(map[int]int)
	for _, val := range rankingData.Data.Rank {
		lxMap[val.Uinfo.Id] += val.RightNum
	}
	roomSumRnum := 0
	for uid, rnum := range lxMap {
		intoLx := user_asset.UserAsset{
			Uid: uid,
			Lx:  rnum * 10,
		}
		store.AddLxChan <- intoLx
		roomSumRnum += rnum
	}
	roomUid, _ := strconv.Atoi(s.Room)
	intoLx := user_asset.UserAsset{
		Uid: roomUid,
		Lx:  roomSumRnum * 10,
	}
	store.AddLxChan <- intoLx
	ctx, cancel := context.WithCancel(context.Background())
	num := s.countdowner(ctx, RANKTIME) //排行榜展示30s
	for n := range num {
		if n < 0 {
			cancel()
			break
		}
		rankingData.Data.Sec = n
		s.BroadcastChan <- tools.FastFunc.DataToBytes(rankingData)
	}

	store.ROOM_READY_STATUS.Delete(s.Room)

	store.ROOM_TIMER_CANCEL.Delete(s.Room)
	store.ROOM_CONTENT_DATA.Delete(s.Room)
	store.ROOM_USER_ANSWER_DATA.Delete(s.Room)
	store.ROOM_COUNTDOWN.Delete(s.Room)
	store.ROOM_RANK_DATA.Delete(s.Room)

	wsBaseMsg2 := EVENT_RESTAGME_DATA{}
	wsBaseMsg2.Event = EVENT_RESTAGME
	s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg2)
	s.noticeRoomStatus()

}

//游戏内循环发题
func (s *Ws) loopGame() {
TAG:
	wsBaseMsg := GAME_READY_START_DATA{}
	wsBaseMsg.Event = EVENT_GAME_START
	sQuestion := s.questionEmiter(s.Room)
	if sQuestion.Id <= 0 { //本轮答题结束
		s.gameEnd()
		return
	}
	store.ROOM_ANSWER_MAP.Store(s.Room, sQuestion)
	ctx, cancel := context.WithCancel(context.Background())
	eatNum := s.countdowner(ctx, 3)
	for n := range eatNum {
		if n < 0 {
			cancel()
			break
		}
		store.ROOM_CONTENT_DATA.Delete(s.Room)
		wsBaseMsg.Data.CountdownSecond = n

		mems := s.getOnlineRoomMems()
		for _, uid := range mems {
			if conn, ok := store.UID_WSCONN_MAP.Load(uid); ok {
				if fmt.Sprintf("%d", uid) == s.Room {
					wsBaseMsg.Data.Question = sQuestion.Question
					wsBaseMsg.Data.QuestionTips = sQuestion.QuestionTips
				} else {
					wsBaseMsg.Data.Question = ""
				}
				conn.(*Ws).Mutex.Lock()
				_ = conn.(*Ws).Conn.WriteMessage(websocket.TextMessage, tools.FastFunc.DataToBytes(wsBaseMsg))
				conn.(*Ws).Mutex.Unlock()
			}
		}
	}
	ctx2, cancel2 := context.WithCancel(context.Background())
	store.ROOM_TIMER_CANCEL.Store(s.Room, cancel2)
	eatNum = s.countdowner(ctx2, ANSWERTIME)
	for n := range eatNum {
		store.ROOM_COUNTDOWN.Store(s.Room, ANSWERTIME-n)
		if n < 0 {
			store.ROOM_ANSWER_MAP.Delete(s.Room)
			store.ROOM_TIMER_CANCEL.Delete(s.Room)
			cancel2()
			s.CreateUserAnswerRecord(0, sQuestion.Id, ANSWERTIME, 0)
			goto TAG
			//ctx3, cancel3 := context.WithCancel(context.Background())
			//restTime := s.countdowner(ctx3, 5)
			//for rt := range restTime {
			//	if rt < 0 {
			//		cancel3()
			//		goto TAG
			//	}
			//	wsBaseMsg.Data.Status = 2 //所有人都没答对，跳题休息时间
			//	wsBaseMsg.Data.CountdownSecond = rt
			//	wsBaseMsg.Data.Question = sQuestion.Question
			//	s.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
			//}
		}
		wsBaseMsg.Data.Status = 1
		wsBaseMsg.Data.CountdownSecond = n

		mems := s.getOnlineRoomMems()
		for _, uid := range mems {
			if conn, ok := store.UID_WSCONN_MAP.Load(uid); ok {
				store.ConsumePowerChan <- int(uid)
				if fmt.Sprintf("%d", uid) == s.Room {
					wsBaseMsg.Data.Question = sQuestion.Question
					wsBaseMsg.Data.QuestionTips = sQuestion.QuestionTips
				} else {
					wsBaseMsg.Data.Question = ""
				}
				conn.(*Ws).Mutex.Lock()
				_ = conn.(*Ws).Conn.WriteMessage(websocket.TextMessage, tools.FastFunc.DataToBytes(wsBaseMsg))
				conn.(*Ws).Mutex.Unlock()
			}
		}
	}
}

func (s *Ws) CreateUserAnswerRecord(uid int, qid int, answerTime int, isRight int) {
	if v, ok := store.ROOM_USER_ANSWER_DATA.Load(s.Room); ok {
		pathData := make([]store.EVENT_DRAW_CONTENT_DATA_DATA, 0)
		gdrawUserAnswers := make([]user_answer.UserAnswer, 0)
		if v1, ok1 := store.ROOM_CONTENT_DATA.Load(s.Room); ok1 {
			pathData = v1.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
		}
		gdrawUserAnswers = v.([]user_answer.UserAnswer)
		gdrawUserAnswers = append(gdrawUserAnswers, user_answer.UserAnswer{
			Uid:        uid,
			AnswerJson: tools.FastFunc.DataToString(pathData),
			QuestionId: qid,
			RoomId:     s.Room,
			AnswerTime: answerTime,
			IsRight:    isRight,
		})
		store.ROOM_USER_ANSWER_DATA.Store(s.Room, gdrawUserAnswers)
		return
	}
	pathData := make([]store.EVENT_DRAW_CONTENT_DATA_DATA, 0)
	if v1, ok1 := store.ROOM_CONTENT_DATA.Load(s.Room); ok1 {
		pathData = v1.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
	}
	gdrawUserAnswers := make([]user_answer.UserAnswer, 0)
	gdrawUserAnswers = append(gdrawUserAnswers, user_answer.UserAnswer{
		Uid:        uid,
		AnswerJson: tools.FastFunc.DataToString(pathData),
		QuestionId: qid,
		RoomId:     s.Room,
		AnswerTime: answerTime,
		IsRight:    isRight,
	})
	store.ROOM_USER_ANSWER_DATA.Store(s.Room, gdrawUserAnswers)
	return
}

func (s *Ws) getRanking() (rankData models.RankData) {
	rankData.Qnum = QUESTIONNUM
	if v, ok := store.ROOM_RANK_DATA.Load(s.Room); ok {
		return v.(models.RankData)
	}
	unUid := make(map[int]int)
	if v, ok := store.ROOM_USER_ANSWER_DATA.Load(s.Room); ok {
		gdrawUserAnswers := v.([]user_answer.UserAnswer)
		for _, val := range gdrawUserAnswers {
			store.GdrawDataChan <- val
			if _, geted := unUid[val.Uid]; geted {
				continue
			}
			if val.Uid == 0 {
				continue
			}
			unUid[val.Uid] = 1
			rightNum := s.getUserAnswerRightNum(val.Uid)
			item := models.UserGameData{}
			item.RightNum = rightNum
			item.Uinfo, _ = user.NewModel(s.mysql.Get()).GetWithId(val.Uid)
			rankData.Rank = append(rankData.Rank, item) //前端排序
		}
		store.ROOM_RANK_DATA.Store(s.Room, rankData)
	}
	return
}

func (s *Ws) getUserAnswerRightNum(uid int) (runm int) {
	userAnswers := make([]user_answer.UserAnswer, 0)
	if v, ok := store.ROOM_USER_ANSWER_DATA.Load(s.Room); ok {
		userAnswers = v.([]user_answer.UserAnswer)
		for _, val := range userAnswers {
			if val.Uid == uid && val.IsRight == 1 {
				runm++
			}
		}
	}
	return
}
