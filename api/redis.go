package api

import (
	"fmt"
	"github.com/go-redis/redis"
	"server_gdraw/configs"
	"sync"
)

var Rds rdser

func init() {
	Rds = &rds{
		roomIdKey:      "room_id",
		QuestionIdsKey: "question_ids",
		RoomUidRoleKey: "room_uid_role",
	}
}

type rdser interface {
	Get() *redis.Client
	//GetRoomIdKey(uid int64) (k string)
	GetQuestionIdsKey() (k string)
	GetRoomUidRoleKey(roomId string, uid int64) (k string)
}

type rds struct {
	sync.Once
	client         *redis.Client
	roomIdKey      string
	QuestionIdsKey string
	RoomUidRoleKey string
}

func (this *rds) Get() *redis.Client {
	this.Do(func() {
		this.client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", configs.Dft.Get().Redis.Host, configs.Dft.Get().Redis.Port),
			Password: configs.Dft.Get().Redis.Pass,
			DB:       0,
		})
	})
	return this.client
}

//func (this *rds) GetRoomIdKey(uid int64) (k string) {
//	k = fmt.Sprintf("%s:%d", this.roomIdKey, uid)
//	return
//}

func (this *rds) GetQuestionIdsKey() (k string) {
	k = fmt.Sprintf("%s", this.QuestionIdsKey)
	return
}

func (this *rds) GetRoomUidRoleKey(roomId string, uid int64) (k string) {
	k = fmt.Sprintf("%s:%s:%d", this.RoomUidRoleKey, roomId, uid)
	return
}
