package services

import (
	"errors"
	"server_gdraw/internal/models"
	"server_gdraw/pkg/tools/store"
	"strconv"
)

var GdrawRoomer gdrawRoomer

func init() {
	GdrawRoomer = &gdrawRoom{}
}

type gdrawRoom struct {
}

type gdrawRoomer interface {
	GetRoomList() (res GetRoomListRes)
	JoinRoom(uid int64, roomId string) (err error)
}

type GetRoomListRes struct {
	List []GetRoomListItem `json:"list"`
}

type GetRoomListItem struct {
	Uinfo      models.GdrawUser       `json:"uinfo"`
	RoomInfo   EVENT_GAME_STATUS_DATA `json:"room_info"`
	OnlinePnum int                    `json:"online_pnum"`
}

func (this *gdrawRoom) GetRoomList() (res GetRoomListRes) { //肯定房主才会来创建房间
	store.ROOM_ALL_PERSON.Range(func(key, value interface{}) bool {
		roomUidInt, _ := strconv.Atoi(key.(string))
		roomLeaderInfos, _ := models.GdrawUserer.FindGdrawUserWithId(roomUidInt)
		if _, ok := store.UID_WSCONN_MAP.Load(int64(roomUidInt)); ok && roomLeaderInfos.Id != 0 {
			item := GetRoomListItem{}
			item.Uinfo = roomLeaderInfos
			if v, oked := store.ROOM_READY_STATUS.Load(key.(string)); oked {
				item.RoomInfo = v.(EVENT_GAME_STATUS_DATA)
			}
			item.OnlinePnum = len(value.([]int64))
			res.List = append(res.List, item)
		}
		return true
	})
	return
}

func (this *gdrawRoom) JoinRoom(uid int64, roomId string) (err error) {
	roomLeaderUid, _ := strconv.Atoi(roomId)
	roomLeaderInfo, _ := models.GdrawUserer.FindGdrawUserWithId(roomLeaderUid)
	if roomLeaderInfo.Id == 0 {
		err = errors.New("房间已不存在")
		return
	}
	var uids []int64
	if v, ok := store.ROOM_ALL_PERSON.Load(roomId); ok {
		uids = v.([]int64)
		for _, item := range uids {
			if item == uid {
				return
			}
		}
		uids = append(uids, uid)
		store.ROOM_ALL_PERSON.Store(roomId, uids)
		return
	}
	uids = append(uids, uid)
	store.ROOM_ALL_PERSON.Store(roomId, uids)
	return
}
