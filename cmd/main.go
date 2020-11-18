package main

import (
	"fmt"
	"server_gdraw/internal/routers"
	"server_gdraw/store"
	"time"
)

func init() {
	redis := new(store.Redis)
	err := redis.Get().Ping().Err()
	if err != nil {
		panic(err)
	}
	mysql := new(store.Mysql)
	mysql.Get()
	log := new(store.Log)
	log.Get().Debug("server started at %v", time.Now())
}

func main() {
	config := new(store.Config)
	err := routers.Init().Run(fmt.Sprintf(":%d", config.Get().Http.Port))
	if err != nil {
		panic(err)
	}
}
