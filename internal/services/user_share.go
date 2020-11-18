package services

import (
	"server_gdraw/internal/models/user_asset"
	"server_gdraw/internal/models/user_share"
)

type UserShare struct {
	base
}

func (s *UserShare) Create(in user_share.UserShare) (res user_share.UserShare, err error) {
	in.Power = 3600
	_, err = user_share.NewModel(s.mysql.Get()).Create(in)
	if err == nil {
		_ = user_asset.NewModel(s.mysql.Get()).AddPowerWithUid(in.Uid, in.Power)
	}
	return
}
