package services

import "server_gdraw/internal/models"

var GdrawIndexer gdrawIndexer

func init() {
	GdrawIndexer = &gdrawIndex{}
}

type gdrawIndex struct {
	AppId     string
	Secret    string
	JwtSecret string
}

type gdrawIndexer interface {
	GetIndexData(uid int) (res GetIndexDataRes, err error)
}

type GetIndexDataRes struct {
	UserInfo  models.GdrawUser      `json:"user_info"`
	AssetInfo models.GdrawUserAsset `json:"asset_info"`
}

func (this *gdrawIndex) GetIndexData(uid int) (res GetIndexDataRes, err error) {
	res.AssetInfo, err = models.GdrawUserAsseter.GetUserAsset(uid)
	if err != nil {
		return
	}
	res.UserInfo, err = models.GdrawUserer.FindGdrawUserWithId(uid)
	if err != nil {
		return
	}
	return
}
