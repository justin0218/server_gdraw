package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"server_gdraw/internal/models"
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
	Event int              `json:"event"`
	Uinfo models.GdrawUser `json:"uinfo"`
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
	Event int              `json:"event"`
	Data  models.GdrawUser `json:"data"`
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
	Role     int              `json:"role"`
	UserInfo models.GdrawUser `json:"user_info"`
}

type gdrawWsBaseMsg struct {
	Event int `json:"event"`
}

type gdrawWs struct {
	Conn                 *websocket.Conn
	Room                 string
	Uid                  int64
	MsgChan              chan []byte
	Token                string
	ErrMsgChan           chan []byte
	BroadcastChan        chan []byte
	BroadcastOutSelfChan chan []byte
	Mutex                sync.Mutex
}

func NewGdrawWs(conn *websocket.Conn, room string, uid int64, token string) GdrawWser {
	return &gdrawWs{
		Conn:                 conn,
		Room:                 room,
		Uid:                  uid,
		MsgChan:              make(chan []byte),
		ErrMsgChan:           make(chan []byte),
		BroadcastChan:        make(chan []byte),
		BroadcastOutSelfChan: make(chan []byte),
		Token:                token,
	}
}

type GdrawWser interface {
	Read()
	Write()
}

//单独携程写消息
func (this *gdrawWs) Write() {
	defer func() {
		_ = this.Conn.Close()
		//this.Conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(3))
	}()
OUT:
	for {
		select {
		case message, ok := <-this.MsgChan:
			if !ok {
				break OUT
			}
			this.Mutex.Lock()
			_ = this.Conn.WriteMessage(websocket.TextMessage, message)
			this.Mutex.Unlock()
			break
		case message, ok := <-this.ErrMsgChan:
			if !ok {
				break OUT
			}
			if conn, ok := store.UID_WSCONN_MAP.Load(this.Uid); ok {
				conn.(*gdrawWs).Mutex.Lock()
				_ = conn.(*gdrawWs).Conn.WriteMessage(websocket.TextMessage, message)
				conn.(*gdrawWs).Mutex.Unlock()
			}
			break OUT
		case message, ok := <-this.BroadcastChan:
			if !ok {
				break OUT
			}
			members := this.getOnlineRoomMems()
			for _, uidInt64 := range members {
				if conn, ok := store.UID_WSCONN_MAP.Load(uidInt64); ok {
					conn.(*gdrawWs).Mutex.Lock()
					_ = conn.(*gdrawWs).Conn.WriteMessage(websocket.TextMessage, message)
					conn.(*gdrawWs).Mutex.Unlock()
				}
			}
			break
		case message, ok := <-this.BroadcastOutSelfChan:
			if !ok {
				break OUT
			}
			members := this.getOnlineRoomMems()
			for _, uid := range members {
				if uid != this.Uid {
					if conn, ok := store.UID_WSCONN_MAP.Load(uid); ok {
						conn.(*gdrawWs).MsgChan <- message
					}
				}
			}
		}
	}
}

//单独携程读取消息
func (this *gdrawWs) Read() {
	this.inroom()
	defer this.outroom()
	msgchan := make(chan []byte)
	for {
		go func() {
			select {
			case m := <-msgchan:
				wsBaseMsg := gdrawWsBaseMsg{}
				err := json.Unmarshal(m, &wsBaseMsg)
				if err != nil {
					this.ErrMsgChan <- []byte(fmt.Sprintf("参数错误:%v", err))
					return
				}
				switch wsBaseMsg.Event {
				case EVENT_DRAW_CONTENT:
					if this.Room != fmt.Sprintf("%d", this.Uid) {
						this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("房主才有权限画画哦"))
						return
					}

					if v, ok := store.ROOM_READY_STATUS.Load(this.Room); ok {
						if v.(EVENT_GAME_STATUS_DATA).RoomStatus != GAME_STARTED {
							this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前游戏状态不能画画哦"))
							return
						}
					} else {
						this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前游戏状态不能画画哦"))
						return
					}

					wsBaseMsg := store.EVENT_DRAW_CONTENT_DATA{}
					_ = json.Unmarshal(m, &wsBaseMsg)
					wsBaseMsg.Event = EVENT_DRAW_CONTENT
					if pathData, ok := store.ROOM_CONTENT_DATA.Load(this.Room); ok {
						pdata := pathData.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
						pdata = append(pdata, wsBaseMsg.Data)
						store.ROOM_CONTENT_DATA.Store(this.Room, pdata)
					} else {
						pdata := []store.EVENT_DRAW_CONTENT_DATA_DATA{}
						pdata = append(pdata, wsBaseMsg.Data)
						store.ROOM_CONTENT_DATA.Store(this.Room, pdata)
					}
					this.BroadcastOutSelfChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
					break
				case EVENT_GAME_START:
					if this.Room != fmt.Sprintf("%d", this.Uid) {
						this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("房主才有权限开始游戏"))
						return
					}
					if v, geted := store.ROOM_READY_STATUS.Load(this.Room); geted {
						if v.(EVENT_GAME_STATUS_DATA).RoomStatus != GAME_WAITING_START && v.(EVENT_GAME_STATUS_DATA).ReadyStatus != GAME_START_STATUS_CANSTART {
							this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前状态不能开始游戏"))
							return
						}
					} else {
						this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("当前状态不能开始游戏"))
						return
					}

					//models.GdrawAnswerLocker.InitLock(this.Room)

					wsBaseMsg := EVENT_GAME_STATUS_DATA{}
					wsBaseMsg.Event = EVENT_GAME_STATUS
					wsBaseMsg.ReadyStatus = GAME_START_STATUS_CANSTART
					wsBaseMsg.RoomStatus = GAME_STARTED
					store.ROOM_READY_STATUS.Store(this.Room, wsBaseMsg)
					this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)

					//选题入库
					randQs, _ := GdrawQuestioner.GetRandQuestions()
					for _, q := range randQs {
						_, _ = models.GdrawRoomQuestioner.Create(models.GdrawRoomQuestion{
							RoomId:       this.Room,
							Question:     q.Question,
							QuestionId:   q.Id,
							QuestionTips: q.QuestionTips,
						})
					}
					this.loopGame()
					break
				case EVENT_GAME_ANSWER:
					if this.Room == fmt.Sprintf("%d", this.Uid) {
						this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("房主不能回答问题哦"))
						return
					}
					wsBaseMsg := EVENT_GAME_ANSWER_DATA{}
					_ = json.Unmarshal(m, &wsBaseMsg)

					tempCon := wsBaseMsg.Data.Answer
					tempCon = strings.ReplaceAll(tempCon, " ", "")
					tempCon = strings.ReplaceAll(tempCon, "\t", "")
					tempCon = strings.ReplaceAll(tempCon, "\n", "")
					if tempCon == "" {
						this.noticeSelfErr(EVENT_ERR_OTHER, fmt.Sprintf("请输入合法内容"))
						return
					}
					if v, ok := store.ROOM_ANSWER_MAP.Load(this.Room); ok {
						if v.(models.GdrawQuestion).Question == wsBaseMsg.Data.Answer {
							//locker := models.GdrawAnswerLocker.Lock(this.Room)
							//if !locker {
							//	return
							//}
							//models.GdrawAnswerLocker.UnLock(this.Room)
							store.ROOM_ANSWER_MAP.Delete(this.Room) //删除答案
							if cancel, ok := store.ROOM_TIMER_CANCEL.Load(this.Room); ok {
								cancel.(context.CancelFunc)()
								if glCountDown, oked := store.ROOM_COUNTDOWN.Load(this.Room); oked {
									this.CreateUserAnswerRecord(wsBaseMsg.Data.Uid, v.(models.GdrawQuestion).Id, glCountDown.(int), 1)
								}
								wsBaseMsg.Event = EVENT_ANSWERRIGHT
								wsBaseMsg.Data.AnswerResult = 1
								//展示答对的信息
								cx, cl := context.WithCancel(context.Background())
								showTime := this.countdowner(cx, 4)
								for s := range showTime {
									if s < 0 {
										cl()
										break
									}
									wsBaseMsg.Data.Sec = s
									this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
								}
								this.loopGame()
							}
						} else {
							this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
						}
					} else {
						this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
					}
					break
				case EVENT_CLEAR_CANVAS:
					wsBaseMsg := EVENT_CLEAR_CANVAS_DATA{}
					wsBaseMsg.Event = EVENT_CLEAR_CANVAS
					store.ROOM_CONTENT_DATA.Delete(this.Room)
					this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
					break
				case EVENT_TIPS:
					this.BroadcastChan <- m
					break

				}
			}
		}()
		_, rmsg, err := this.Conn.ReadMessage()
		if err != nil {
			break
		}
		msgchan <- rmsg
	}
}

//通知历史的路径数据
func (this *gdrawWs) noticeRoomPathData() {
	wsBaseMsg := store.EVENT_INITPATH_DATA_DATA{}
	wsBaseMsg.Event = EVENT_INITPATH_DATA
	if v, ok := store.ROOM_CONTENT_DATA.Load(this.Room); ok {
		wsBaseMsg.Data = v.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
	}
	this.MsgChan <- []byte(tools.FastFunc.DataToString(wsBaseMsg))
}

//通知错误
func (this *gdrawWs) noticeSelfErr(eventCode int, errMsg string) {
	wsBaseMsg := EVENT_ERR_DATA{}
	wsBaseMsg.Event = eventCode
	wsBaseMsg.Msg = errMsg
	this.MsgChan <- []byte(tools.FastFunc.DataToString(wsBaseMsg))
}

//通知房间内成员
func (this *gdrawWs) noticeMembersToRoom() {
	wsBaseMsg := EVENT_MEMBER_LIST_DATA{}
	wsBaseMsg.Event = EVENT_MEMBER_LIST
	wsBaseMsg.Data = this.GetMembersWithRoom()
	this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
}

//发送房间状态
func (this *gdrawWs) noticeRoomStatus() {
	if v, ok := store.ROOM_READY_STATUS.Load(this.Room); ok {
		if v.(EVENT_GAME_STATUS_DATA).RoomStatus == GAME_END || v.(EVENT_GAME_STATUS_DATA).RoomStatus == GAME_STARTED {
			this.BroadcastChan <- tools.FastFunc.DataToBytes(v.(EVENT_GAME_STATUS_DATA))
			return
		}
		if v.(EVENT_GAME_STATUS_DATA).RoomStatus == GAME_WAITING_START {
			members := this.getOnlineRoomMems()
			wsBaseMsg := EVENT_GAME_STATUS_DATA{}
			wsBaseMsg.Event = EVENT_GAME_STATUS
			if len(members) > 1 {
				//可以开始游戏
				wsBaseMsg.ReadyStatus = GAME_START_STATUS_CANSTART
			} else {
				wsBaseMsg.ReadyStatus = GAME_START_STATUS_NOCANSTART
			}
			wsBaseMsg.RoomStatus = GAME_WAITING_START
			store.ROOM_READY_STATUS.Store(this.Room, wsBaseMsg)
			this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
		}
		return
	}
	members := this.getOnlineRoomMems()
	wsBaseMsg := EVENT_GAME_STATUS_DATA{}
	wsBaseMsg.Event = EVENT_GAME_STATUS
	if len(members) > 1 {
		//可以开始游戏
		wsBaseMsg.ReadyStatus = GAME_START_STATUS_CANSTART
	} else {
		wsBaseMsg.ReadyStatus = GAME_START_STATUS_NOCANSTART
	}
	wsBaseMsg.RoomStatus = GAME_WAITING_START
	store.ROOM_READY_STATUS.Store(this.Room, wsBaseMsg)
	this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
	return
}

//获取房间成员
func (this *gdrawWs) GetMembersWithRoom() (data []EVENT_MEMBER_LIST_DATA_DATA) {
	uids := this.getOnlineRoomMems()
	membersFull, _ := models.GdrawUserer.FindGdrawUserWithIds(uids)
	for _, val := range membersFull {
		item := EVENT_MEMBER_LIST_DATA_DATA{}
		item.UserInfo = val
		data = append(data, item)
	}
	return
}

//进入房间
func (this *gdrawWs) inroom() {
	err := GdrawUserer.VerifyTokenFunc(this.Uid, this.Token)
	if err != nil {
		this.MsgChan <- []byte(fmt.Sprintf(`{"event":-1,"msg":"%s"}`, err.Error()))
		return
	}
	//效验体力
	userAsset, _ := models.GdrawUserAsseter.GetUserAsset(int(this.Uid))
	if userAsset.Power <= 0 {
		this.MsgChan <- []byte(fmt.Sprintf(`{"event":-1,"msg":"体力值不足"}`))
		return
	}
	err = GdrawRoomer.JoinRoom(this.Uid, this.Room)
	if err != nil {
		this.MsgChan <- []byte(fmt.Sprintf(`{"event":-1,"msg":"%s"}`, err.Error()))
		return
	}
	store.UID_WSCONN_MAP.Store(this.Uid, this)

	this.noticeSelfRole()
	//通知房间内所有成员当前房间状态
	this.noticeRoomStatus()
	//进入房间就需要通知的
	this.noticeRoomLeaderInfo()
	this.noticeMembersToRoom()
	this.noticeRoomPathData()

	noticeMsg := EVENT_NOTICE_INOUT_DATA{}
	noticeMsg.Event = EVENT_NOTICE_IN
	noticeMsg.Uinfo, _ = models.GdrawUserer.FindGdrawUserWithId(int(this.Uid))
	this.BroadcastChan <- tools.FastFunc.DataToBytes(noticeMsg)
	return
}

func (this *gdrawWs) noticeSelfRole() {
	wsBaseMsg := EVENT_USER_ROLE_DATA{}
	wsBaseMsg.Event = EVENT_USER_ROLE
	if this.Room == fmt.Sprintf("%d", this.Uid) {
		wsBaseMsg.Role = 1
	}
	this.MsgChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
}

//获取房主的信息
func (this *gdrawWs) noticeRoomLeaderInfo() (uninfo models.GdrawUser) {
	uidInt, _ := strconv.Atoi(this.Room)
	uninfo, _ = models.GdrawUserer.FindGdrawUserWithId(uidInt)
	wsBaseMsg := EVENT_ROOMLEADER_DATA{}
	wsBaseMsg.Event = EVENT_ROOMLEADER_INFO
	wsBaseMsg.Data = uninfo
	this.MsgChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
	return
}

//离开房间
func (this *gdrawWs) outroom() {
	if _, ok := store.UID_WSCONN_MAP.Load(this.Uid); ok {
		store.UID_WSCONN_MAP.Delete(this.Uid)
		uids := this.getOnlineRoomMems()
		for i, uid := range uids {
			if uid == this.Uid {
				uids = append(uids[:i], uids[i+1:]...)
				break
			}
		}
		store.ROOM_ALL_PERSON.Store(this.Room, uids)
		this.noticeMembersToRoom()
		//通知房间内所有成员当前房间状态
		this.noticeRoomStatus()

		//通知离开房间
		noticeMsg := EVENT_NOTICE_INOUT_DATA{}
		noticeMsg.Event = EVENT_NOTICE_OUT
		noticeMsg.Uinfo, _ = models.GdrawUserer.FindGdrawUserWithId(int(this.Uid))
		this.BroadcastChan <- tools.FastFunc.DataToBytes(noticeMsg)
	}
}

func (this *gdrawWs) getOnlineRoomMems() (uids []int64) {
	if v, ok := store.ROOM_ALL_PERSON.Load(this.Room); ok {
		uids = v.([]int64)
		return
	}
	return
}

//游戏结束
func (this *gdrawWs) gameEnd() {
	store.ROOM_ANSWER_MAP.Delete(this.Room)
	wsBaseMsg := EVENT_GAME_STATUS_DATA{}
	wsBaseMsg.Event = EVENT_GAME_STATUS
	wsBaseMsg.RoomStatus = GAME_END
	wsBaseMsg.ReadyStatus = GAME_START_STATUS_NOCANSTART
	store.ROOM_READY_STATUS.Store(this.Room, wsBaseMsg)
	this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)

	rankingData := EVENT_RANKING_DATA{}
	rankingData.Event = EVENT_RANKING
	rankingData.Data = this.getRanking()

	lxMap := make(map[int]int)
	for _, val := range rankingData.Data.Rank {
		lxMap[val.Uinfo.Id] += val.RightNum
	}
	roomSumRnum := 0
	for uid, rnum := range lxMap {
		intoLx := models.GdrawUserAsset{
			Uid: uid,
			Lx:  rnum * 10,
		}
		store.AddLxChan <- intoLx
		roomSumRnum += rnum
	}
	roomUid, _ := strconv.Atoi(this.Room)
	intoLx := models.GdrawUserAsset{
		Uid: roomUid,
		Lx:  roomSumRnum * 10,
	}
	store.AddLxChan <- intoLx
	ctx, cancel := context.WithCancel(context.Background())
	num := this.countdowner(ctx, RANKTIME) //排行榜展示30s
	for n := range num {
		if n < 0 {
			cancel()
			break
		}
		rankingData.Data.Sec = n
		this.BroadcastChan <- tools.FastFunc.DataToBytes(rankingData)
	}

	store.ROOM_READY_STATUS.Delete(this.Room)

	store.ROOM_TIMER_CANCEL.Delete(this.Room)
	store.ROOM_CONTENT_DATA.Delete(this.Room)
	store.ROOM_USER_ANSWER_DATA.Delete(this.Room)
	store.ROOM_COUNTDOWN.Delete(this.Room)
	store.ROOM_RANK_DATA.Delete(this.Room)

	wsBaseMsg2 := EVENT_RESTAGME_DATA{}
	wsBaseMsg2.Event = EVENT_RESTAGME
	this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg2)
	this.noticeRoomStatus()

}

//游戏内循环发题
func (this *gdrawWs) loopGame() {
TAG:
	wsBaseMsg := GAME_READY_START_DATA{}
	wsBaseMsg.Event = EVENT_GAME_START
	thisQuestion := this.questionEmiter(this.Room)
	if thisQuestion.Id <= 0 { //本轮答题结束
		this.gameEnd()
		return
	}
	store.ROOM_ANSWER_MAP.Store(this.Room, thisQuestion)
	ctx, cancel := context.WithCancel(context.Background())
	eatNum := this.countdowner(ctx, 3)
	for n := range eatNum {
		if n < 0 {
			cancel()
			break
		}
		store.ROOM_CONTENT_DATA.Delete(this.Room)
		wsBaseMsg.Data.CountdownSecond = n

		mems := this.getOnlineRoomMems()
		for _, uid := range mems {
			if conn, ok := store.UID_WSCONN_MAP.Load(uid); ok {
				if fmt.Sprintf("%d", uid) == this.Room {
					wsBaseMsg.Data.Question = thisQuestion.Question
					wsBaseMsg.Data.QuestionTips = thisQuestion.QuestionTips
				} else {
					wsBaseMsg.Data.Question = ""
				}
				conn.(*gdrawWs).Mutex.Lock()
				_ = conn.(*gdrawWs).Conn.WriteMessage(websocket.TextMessage, tools.FastFunc.DataToBytes(wsBaseMsg))
				conn.(*gdrawWs).Mutex.Unlock()
			}
		}
	}
	ctx2, cancel2 := context.WithCancel(context.Background())
	store.ROOM_TIMER_CANCEL.Store(this.Room, cancel2)
	eatNum = this.countdowner(ctx2, ANSWERTIME)
	for n := range eatNum {
		store.ROOM_COUNTDOWN.Store(this.Room, ANSWERTIME-n)
		if n < 0 {
			store.ROOM_ANSWER_MAP.Delete(this.Room)
			store.ROOM_TIMER_CANCEL.Delete(this.Room)
			cancel2()
			this.CreateUserAnswerRecord(0, thisQuestion.Id, ANSWERTIME, 0)
			goto TAG
			//ctx3, cancel3 := context.WithCancel(context.Background())
			//restTime := this.countdowner(ctx3, 5)
			//for rt := range restTime {
			//	if rt < 0 {
			//		cancel3()
			//		goto TAG
			//	}
			//	wsBaseMsg.Data.Status = 2 //所有人都没答对，跳题休息时间
			//	wsBaseMsg.Data.CountdownSecond = rt
			//	wsBaseMsg.Data.Question = thisQuestion.Question
			//	this.BroadcastChan <- tools.FastFunc.DataToBytes(wsBaseMsg)
			//}
		}
		wsBaseMsg.Data.Status = 1
		wsBaseMsg.Data.CountdownSecond = n

		mems := this.getOnlineRoomMems()
		for _, uid := range mems {
			if conn, ok := store.UID_WSCONN_MAP.Load(uid); ok {
				store.ConsumePowerChan <- int(uid)
				if fmt.Sprintf("%d", uid) == this.Room {
					wsBaseMsg.Data.Question = thisQuestion.Question
					wsBaseMsg.Data.QuestionTips = thisQuestion.QuestionTips
				} else {
					wsBaseMsg.Data.Question = ""
				}
				conn.(*gdrawWs).Mutex.Lock()
				_ = conn.(*gdrawWs).Conn.WriteMessage(websocket.TextMessage, tools.FastFunc.DataToBytes(wsBaseMsg))
				conn.(*gdrawWs).Mutex.Unlock()
			}
		}
	}
}

func (this *gdrawWs) CreateUserAnswerRecord(uid int, qid int, answerTime int, isRight int) {
	if v, ok := store.ROOM_USER_ANSWER_DATA.Load(this.Room); ok {
		pathData := []store.EVENT_DRAW_CONTENT_DATA_DATA{}
		gdrawUserAnswers := []models.GdrawUserAnswer{}
		if v1, ok1 := store.ROOM_CONTENT_DATA.Load(this.Room); ok1 {
			pathData = v1.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
		}
		gdrawUserAnswers = v.([]models.GdrawUserAnswer)
		gdrawUserAnswers = append(gdrawUserAnswers, models.GdrawUserAnswer{
			Uid:        uid,
			AnswerJson: tools.FastFunc.DataToString(pathData),
			QuestionId: qid,
			RoomId:     this.Room,
			AnswerTime: answerTime,
			IsRight:    isRight,
		})
		store.ROOM_USER_ANSWER_DATA.Store(this.Room, gdrawUserAnswers)
		return
	}
	pathData := []store.EVENT_DRAW_CONTENT_DATA_DATA{}
	if v1, ok1 := store.ROOM_CONTENT_DATA.Load(this.Room); ok1 {
		pathData = v1.([]store.EVENT_DRAW_CONTENT_DATA_DATA)
	}
	gdrawUserAnswers := []models.GdrawUserAnswer{}
	gdrawUserAnswers = append(gdrawUserAnswers, models.GdrawUserAnswer{
		Uid:        uid,
		AnswerJson: tools.FastFunc.DataToString(pathData),
		QuestionId: qid,
		RoomId:     this.Room,
		AnswerTime: answerTime,
		IsRight:    isRight,
	})
	store.ROOM_USER_ANSWER_DATA.Store(this.Room, gdrawUserAnswers)
	return
}

func (this *gdrawWs) getRanking() (rankData models.RankData) {
	rankData.Qnum = QUESTIONNUM
	if v, ok := store.ROOM_RANK_DATA.Load(this.Room); ok {
		return v.(models.RankData)
	}
	unUid := make(map[int]int)
	if v, ok := store.ROOM_USER_ANSWER_DATA.Load(this.Room); ok {
		gdrawUserAnswers := v.([]models.GdrawUserAnswer)
		for _, val := range gdrawUserAnswers {
			store.GdrawDataChan <- val
			if _, geted := unUid[val.Uid]; geted {
				continue
			}
			if val.Uid == 0 {
				continue
			}
			unUid[val.Uid] = 1
			rightNum := this.getUserAnswerRightNum(val.Uid)
			item := models.UserGameData{}
			item.RightNum = rightNum
			item.Uinfo, _ = models.GdrawUserer.FindGdrawUserWithId(val.Uid)
			rankData.Rank = append(rankData.Rank, item) //前端排序
		}
		store.ROOM_RANK_DATA.Store(this.Room, rankData)
	}
	return
}

func (this *gdrawWs) getUserAnswerRightNum(uid int) (runm int) {
	userAnswers := []models.GdrawUserAnswer{}
	if v, ok := store.ROOM_USER_ANSWER_DATA.Load(this.Room); ok {
		userAnswers = v.([]models.GdrawUserAnswer)
		for _, val := range userAnswers {
			if val.Uid == uid && val.IsRight == 1 {
				runm++
			}
		}
	}
	return
}
