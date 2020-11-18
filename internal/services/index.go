package services

import (
	"server_gdraw/internal/models/user"
	"server_gdraw/internal/models/user_asset"
)

type Index struct {
	base
}

type GetIndexDataRes struct {
	UserInfo  user.User            `json:"user_info"`
	AssetInfo user_asset.UserAsset `json:"asset_info"`
}

func (s *Index) GetIndexData(uid int) (res GetIndexDataRes, err error) {
	db := s.mysql.Get()
	res.AssetInfo, err = user_asset.NewModel(db).GetOrCreateWithUid(uid)
	if err != nil {
		return
	}
	res.UserInfo, err = user.NewModel(db).GetWithId(uid)
	if err != nil {
		return
	}
	return
}
