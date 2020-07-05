package models

type UserGameData struct {
	RightNum int       `json:"right_num"`
	Role     int       `json:"role"`
	Uinfo    GdrawUser `json:"uinfo"`
}

type RankData struct {
	Qnum int            `json:"qnum"`
	Sec  int            `json:"sec"`
	Rank []UserGameData `json:"rank"`
}
