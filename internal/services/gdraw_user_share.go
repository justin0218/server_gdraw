package services

import (
	"server_gdraw/internal/models"
)

var GdrawUserShareer gdrawUserShareer

func init() {
	GdrawUserShareer = &gdrawUserShare{
		TableName: "gdraw_user_shares",
	}
}

type gdrawUserShare struct {
	TableName string
}

type gdrawUserShareer interface {
	Create(in models.GdrawUserShare) (res models.GdrawUserShare, err error)
}

func (this *gdrawUserShare) Create(in models.GdrawUserShare) (res models.GdrawUserShare, err error) {
	in.Power = 3600
	_, err = models.GdrawUserShareer.Create(in)
	if err == nil {
		_ = models.GdrawUserAsseter.AddUserPower(in.Uid, in.Power)
	}
	return
}
