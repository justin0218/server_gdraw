package services

import "server_gdraw/internal/models"

var GdrawFeedbacker gdrawFeedbacker

func init() {
	GdrawFeedbacker = &gdrawFeedback{}
}

type gdrawFeedback struct {
	AppId     string
	Secret    string
	JwtSecret string
}

type gdrawFeedbacker interface {
	Create(in models.GdrawFeedback) (res models.GdrawFeedback, err error)
}

func (this *gdrawFeedback) Create(in models.GdrawFeedback) (res models.GdrawFeedback, err error) {
	return models.GdrawFeedbacker.Create(in)
}
