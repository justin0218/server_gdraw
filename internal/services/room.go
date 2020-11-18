package services

import (
	"errors"
	"server_gdraw/internal/models/user"
	"server_gdraw/pkg/tools/store"
	"strconv"
)

type Room struct {
	base
}

type GetRoomListRes struct {
	List []GetRoomListItem `json:"list"`
}

type GetRoomListItem struct {
	Uinfo      user.User              `json:"uinfo"`
	RoomInfo   EVENT_GAME_STATUS_DATA `json:"room_info"`
	OnlinePnum int                    `json:"online_pnum"`
}

func (s *Room) GetRoomList() (res GetRoomListRes) { //肯定房主才会来创建房间
	store.ROOM_ALL_PERSON.Range(func(key, value interface{}) bool {
		roomUidInt, _ := strconv.Atoi(key.(string))
		roomLeaderInfos, _ := user.NewModel(s.mysql.Get()).GetWithId(roomUidInt)
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

func (s *Room) JoinRoom(uid int64, roomId string) (err error) {
	roomLeaderUid, _ := strconv.Atoi(roomId)
	roomLeaderInfo, _ := user.NewModel(s.mysql.Get()).GetWithId(roomLeaderUid)
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
