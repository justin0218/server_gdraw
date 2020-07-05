package services

import (
	"fmt"
	"github.com/gorilla/websocket"
	"server_gdraw/api"
	"server_gdraw/internal/models"
	"server_gdraw/pkg/tools/store"
)

func NewJob() jober {
	return &job{}
}

type job struct {
}

type jober interface {
	GdrawData()
}

func (this *job) GdrawData() {
	for {
		select {
		case msg, ok := <-store.GdrawDataChan:
			if !ok {
				return
			}
			err := models.GdrawUserAnswerer.Create(msg)
			if err != nil {
				api.Log.Get().Error("job GdrawData GdrawUserAnswerer Create err:%v", err)
			}
		case uid, ok := <-store.ConsumePowerChan:
			if !ok {
				return
			}
			err := models.GdrawUserAsseter.UpdateUserPower(uid, 1)
			if err != nil {
				if conn, ok := store.UID_WSCONN_MAP.Load(int64(uid)); ok {
					conn.(*gdrawWs).Mutex.Lock()
					_ = conn.(*gdrawWs).Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"event":-1,"msg":"体力值不足"}`)))
					conn.(*gdrawWs).Mutex.Unlock()
				}
				api.Log.Get().Error("job GdrawData UpdateUserPower err:%v", err)
			}
		case asset, ok := <-store.AddLxChan:
			if !ok {
				return
			}
			err := models.GdrawUserAsseter.UpdateUserLx(asset.Uid, asset.Lx)
			if err != nil {
				api.Log.Get().Error("job GdrawData UpdateUserLx err:%v", err)
			}

		}
	}
}
