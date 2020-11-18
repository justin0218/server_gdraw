package services

import (
	"server_gdraw/internal/models/user"
	"server_gdraw/internal/models/user_asset"
)

var Praises = []string{"你心地真善良", "你真聪明", "你真厉害", "你好棒", "你太帅了", "你酷毙了", "你的思维太活跃了", "你每天都这么精神", "你的人格魅力真强", "你决不是一般人"}

type Rank struct {
	base
}

type GetRankDataRes struct {
	UserInfo  user.User            `json:"user_info"`
	AssetInfo user_asset.UserAsset `json:"asset_info"`
	Praise    string               `json:"praise"`
}

func (s *Rank) GetRankData() (res []GetRankDataRes, err error) {
	var (
		ranks []user_asset.UserAsset
	)
	ranks, err = user_asset.NewModel(s.mysql.Get()).GetWithLimit(10)
	if err != nil {
		return
	}
	for dx, val := range ranks {
		item := GetRankDataRes{}
		item.Praise = Praises[dx]
		item.AssetInfo = val
		item.UserInfo, _ = user.NewModel(s.mysql.Get()).GetWithId(val.Uid)
		res = append(res, item)
	}
	return
}
