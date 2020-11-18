package wxminiapp

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

type CodeRes struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}

func GetAuthorizationCode(code string) (ret CodeRes, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
	_, _, errs := gorequest.New().Get(url).EndStruct(&ret)
	if len(errs) > 0 {
		err = errs[0]
		return
	}
	if ret.Errcode != 0 {
		err = errors.New(ret.ErrMsg)
		return
	}
	return
}
