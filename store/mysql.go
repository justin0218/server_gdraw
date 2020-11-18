package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct {
	conf *Config
}

func (s *Mysql) Get() *gorm.DB {
	mysqlOnce.Do(func() {
		var err error
		mysqlClient, err = gorm.Open("mysql", s.conf.Get().Mysql.Master.Addr)
		if err != nil {
			panic(err)
		}
		if s.conf.Get().Runmode == "debug" {
			mysqlClient.LogMode(true)
		}
		mysqlClient.DB().SetMaxOpenConns(s.conf.Get().Mysql.Master.MaxOpenConns)
		mysqlClient.DB().SetMaxIdleConns(s.conf.Get().Mysql.Master.MaxIdleConns)
	})
	return mysqlClient
}
