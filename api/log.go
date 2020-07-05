package api

import (
	"github.com/astaxie/beego/logs"
	"sync"
)

var Log loger

func init() {
	Log = &log{}
}

type loger interface {
	Get() *logs.BeeLogger
}

type log struct {
	sync.Once
	client *logs.BeeLogger
}

func (this *log) Get() *logs.BeeLogger {
	this.Do(func() {
		logs.SetLevel(0)
		this.client = logs.NewLogger()
		err := this.client.SetLogger(logs.AdapterMultiFile, `{"filename":"./logs/log.log","separate":["error","warning","info","debug"]}`)
		if err != nil {
			panic(err)
		}
		this.client.EnableFuncCallDepth(true)
	})
	return this.client
}
