package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"server_gdraw/internal/models/user"
	"server_gdraw/pkg/jwt"
	"server_gdraw/pkg/wxminiapp"
	"time"
)

type User struct {
	base
}

func (s *User) ParseUserInfo(encryptedData, iv, code string) (userFull user.UserFull, err error) {
	var (
		uinfo      user.User
		wxuserinfo wxminiapp.User
		wxcodeRes  wxminiapp.CodeRes
	)
	wxcodeRes, err = wxminiapp.GetAuthorizationCode(code)
	if err != nil {
		return
	}
	wxuserinfo, err = wxminiapp.DecryptUserInfo(wxcodeRes.SessionKey, encryptedData, iv)
	if err != nil {
		return
	}
	uinfo, err = s.RegisterUser(wxuserinfo)
	if err != nil {
		return
	}
	userFull.Token, err = jwt.CreateToken(int64(uinfo.Id))
	if err != nil {
		return
	}
	userFull.User = uinfo
	return
}

func (s *User) RegisterUser(userinfo wxminiapp.User) (uinfo user.User, err error) {
	if userinfo.Openid == "" {
		userinfo.Openid = fmt.Sprintf("H5-%d", time.Now().UnixNano())
	}
	uinfo, err = user.NewModel(s.mysql.Get()).GetWithOpenid(userinfo.Openid)
	if err == gorm.ErrRecordNotFound {
		inuinfo := user.User{}
		inuinfo.Nickname = userinfo.Nickname
		inuinfo.Openid = userinfo.Openid
		inuinfo.Avatar = userinfo.Avatar
		uinfo, err = user.NewModel(s.mysql.Get()).Create(inuinfo)
		if err != nil {
			return
		}
	}
	return
}
