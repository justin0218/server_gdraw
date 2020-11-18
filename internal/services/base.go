package services

import "server_gdraw/store"

type base struct {
	redis store.Redis
	mysql store.Mysql
}
