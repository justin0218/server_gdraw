package store

import (
	"server_gdraw/internal/models/user_answer"
	"server_gdraw/internal/models/user_asset"
	"sync"
)

var (
	GdrawDataChan    = make(chan user_answer.UserAnswer, 64)
	ConsumePowerChan = make(chan int, 128)
	AddLxChan        = make(chan user_asset.UserAsset, 128)
)

var (
	UID_WSCONN_MAP    = sync.Map{}
	ROOM_READY_STATUS = sync.Map{}
	ROOM_TIMER_CANCEL = sync.Map{}
	ROOM_ANSWER_MAP   = sync.Map{}
	//ROOM_USERROLR_MAP     = sync.Map{}
	ROOM_CONTENT_DATA     = sync.Map{}
	ROOM_USER_ANSWER_DATA = sync.Map{}
	ROOM_COUNTDOWN        = sync.Map{}
	//ROOM_MENBER           = sync.Map{}
	ROOM_RANK_DATA = sync.Map{}
	//ROOM_LEADER_INFO = sync.Map{}
	ROOM_ALL_PERSON = sync.Map{}
)

type EVENT_DRAW_CONTENT_DATA struct { //画笔数据
	Event int                          `json:"event"`
	Data  EVENT_DRAW_CONTENT_DATA_DATA `json:"data"`
}

type EVENT_INITPATH_DATA_DATA struct {
	Event int                            `json:"event"`
	Data  []EVENT_DRAW_CONTENT_DATA_DATA `json:"data"`
}

type EVENT_DRAW_CONTENT_DATA_DATA struct {
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	StrokeStyle  string  `json:"stroke_style"`
	LineWidth    int     `json:"line_width"`
	Action       string  `json:"action"`
	CanvasWidth  int     `json:"canvas_width"`
	CanvasHeight int     `json:"canvas_height"`
}
