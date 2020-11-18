package wxminiapp

import "server_gdraw/store"

var (
	appid  = new(store.Config).Get().Miniapp.Appid
	secret = new(store.Config).Get().Miniapp.Secret
)
