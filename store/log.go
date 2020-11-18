package store

import (
	"github.com/astaxie/beego/logs"
)

type Log struct {
}

func (s *Log) Get() *logs.BeeLogger {
	logOnce.Do(func() {
		logs.SetLevel(0)
		logClient = logs.NewLogger()
		err := logClient.SetLogger(logs.AdapterMultiFile, `{"filename":"./logs/log.log","separate":["error","warning","info","debug"]}`)
		if err != nil {
			panic(err)
		}
		logClient.EnableFuncCallDepth(true)
	})
	return logClient
}
