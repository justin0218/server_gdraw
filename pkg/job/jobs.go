package job

import (
	"fmt"
	"github.com/gorilla/websocket"
	"server_gdraw/api"
	"server_gdraw/internal/models/user_answer"
	"server_gdraw/internal/models/user_asset"
	"server_gdraw/internal/services/gdraw_ws_service"
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
			err := user_answer.GdrawUserAnswerer.Create(msg)
			if err != nil {
				api.Log.Get().Error("job GdrawData GdrawUserAnswerer Create err:%v", err)
			}
		case uid, ok := <-store.ConsumePowerChan:
			if !ok {
				return
			}
			err := user_asset.GdrawUserAsseter.UpdateUserPower(uid, 1)
			if err != nil {
				if conn, ok := store.UID_WSCONN_MAP.Load(int64(uid)); ok {
					conn.(*gdraw_ws_service.gdrawWs).Mutex.Lock()
					_ = conn.(*gdraw_ws_service.gdrawWs).Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"event":-1,"msg":"体力值不足"}`)))
					conn.(*gdraw_ws_service.gdrawWs).Mutex.Unlock()
				}
				api.Log.Get().Error("job GdrawData UpdateUserPower err:%v", err)
			}
		case asset, ok := <-store.AddLxChan:
			if !ok {
				return
			}
			err := user_asset.GdrawUserAsseter.UpdateUserLx(asset.Uid, asset.Lx)
			if err != nil {
				api.Log.Get().Error("job GdrawData UpdateUserLx err:%v", err)
			}

		}
	}
}
