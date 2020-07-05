package services

import "server_gdraw/internal/models"

var GdrawRanker gdrawRanker

func init() {
	GdrawRanker = &gdrawRank{
		Praises: []string{"你心地真善良", "你真聪明", "你真厉害", "你好棒", "你太帅了", "你酷毙了", "你的思bai维太活跃了", "你每天都这么精神", "你的人格魅力真强", "你决不是一般人"},
	}
}

type gdrawRank struct {
	AppId     string
	Secret    string
	JwtSecret string
	Praises   []string
}

type gdrawRanker interface {
	GetRankData() (res []GetRankDataRes, err error)
}

type GetRankDataRes struct {
	UserInfo  models.GdrawUser      `json:"user_info"`
	AssetInfo models.GdrawUserAsset `json:"asset_info"`
	Praise    string                `json:"praise"`
}

func (this *gdrawRank) GetRankData() (res []GetRankDataRes, err error) {
	var (
		ranks []models.GdrawUserAsset
	)
	ranks, err = models.GdrawUserAsseter.GetRankData(10)
	if err != nil {
		return
	}
	for dx, val := range ranks {
		item := GetRankDataRes{}
		item.Praise = this.Praises[dx]
		item.AssetInfo = val
		item.UserInfo, _ = models.GdrawUserer.FindGdrawUserWithId(val.Uid)
		res = append(res, item)
	}
	return
}
