package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"server_gdraw/internal/models"
	"server_gdraw/pkg/resp"
	"server_gdraw/pkg/tools"
	"strconv"
	"time"
)

var GdrawUserer gdrawUserer

func init() {
	GdrawUserer = &gdrawUser{
		AppId:     "wxc7d86da0fc15b41f",
		Secret:    "3ee727775e23c42f872099d605a3f877",
		JwtSecret: "bxy)@)#)&",
	}
}

type gdrawUser struct {
	AppId     string
	Secret    string
	JwtSecret string
}

type gdrawUserer interface {
	ParseUserInfo(encryptedData, iv, code string) (userinfoFull GdrawUserInfoFull, err error)
	VerifyToken() gin.HandlerFunc
	VerifyTokenFunc(uid int64, token string) (err error)
	RegisterUser(userinfo GdrawWxUserInfo) (uinfo models.GdrawUser, err error)
	CreateToken(uid int64) (string, error)
	MakePassword(uid int) (password string, err error)
	Login(uid int, password string) (userinfo models.GdrawUser, loginSuccess bool, err error)
}

type GdrawUserWxCodeRes struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}

type GdrawWxUserInfo struct {
	Openid    string `json:"openid"`
	Nickname  string `json:"nickName"`
	Avatar    string `json:"avatarUrl"`
	UnionId   string `json:"unionId"`
	Watermark struct {
		APPID string `appid`
	} `json:"watermark"`
}

type GdrawUserInfoFull struct {
	Token string `json:"token"`
	models.GdrawUser
}

func (this *gdrawUser) ParseUserInfo(encryptedData, iv, code string) (userinfoFull GdrawUserInfoFull, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", this.AppId, this.Secret, code)
	wxcodeRes := GdrawUserWxCodeRes{}
	_, _, errs := gorequest.New().Get(url).EndStruct(&wxcodeRes)
	if len(errs) != 0 {
		return
	}
	if wxcodeRes.Errcode != 0 {
		err = errors.New(wxcodeRes.ErrMsg)
		return
	}
	var (
		uinfo      models.GdrawUser
		wxuserinfo GdrawWxUserInfo
	)
	wxuserinfo, err = this.DecryptUserInfo(wxcodeRes.SessionKey, encryptedData, iv)
	if err != nil {
		return
	}
	uinfo, err = this.RegisterUser(wxuserinfo)
	if err != nil {
		return
	}
	userinfoFull.Token, err = this.CreateToken(int64(uinfo.Id))
	if err != nil {
		return
	}
	userinfoFull.GdrawUser = uinfo
	return
}

func (this *gdrawUser) RegisterUser(userinfo GdrawWxUserInfo) (uinfo models.GdrawUser, err error) {
	if userinfo.Openid == "" {
		userinfo.Openid = fmt.Sprintf("H5-%d", time.Now().UnixNano())
	}
	uinfo, err = models.GdrawUserer.FindGdrawUserWithOpenid(userinfo.Openid)
	if err == gorm.ErrRecordNotFound {
		inuinfo := models.GdrawUser{}
		inuinfo.Nickname = userinfo.Nickname
		inuinfo.Openid = userinfo.Openid
		inuinfo.Avatar = userinfo.Avatar
		uinfo, err = models.GdrawUserer.CreateGdrawUser(inuinfo)
		if err != nil {
			return
		}
	}
	return
}

func (this *gdrawUser) MakePassword(uid int) (password string, err error) {
	return models.GdrawUserer.MakePassword(uid)
}

func (this *gdrawUser) DecryptUserInfo(sessionKey, encryptedData, iv string) (userinfo GdrawWxUserInfo, err error) {
	if len(sessionKey) != 24 {
		err = errors.New("sessionKey length is error")
		return
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		err = errors.New("DecodeBase64Error")
		return
	}
	if len(iv) != 24 {
		err = errors.New("iv length is error")
		return
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		err = errors.New("DecodeBase64Error" + err.Error())
		return
	}
	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		err = errors.New("DecodeBase64Error" + err.Error())
		return
	}
	aesPlantText := make([]byte, len(aesCipherText))
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		err = errors.New("IllegalBuffer" + err.Error())
		return
	}
	mode := cipher.NewCBCDecrypter(aesBlock, aesIV)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = this.PKCS7UnPadding(aesPlantText)
	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	aesPlantText = []byte(re.ReplaceAllString(string(aesPlantText), "$1"))
	err = json.Unmarshal(aesPlantText, &userinfo)
	if err != nil {
		err = errors.New("DecodeJsonError" + err.Error())
		return
	}
	if userinfo.Watermark.APPID != this.AppId {
		err = errors.New("appID is not match")
		return
	}
	return
}

func (this *gdrawUser) PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}

type CustomClaims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

func (this *gdrawUser) CreateToken(uid int64) (string, error) {
	stringUid := strconv.FormatInt(uid, 10)
	claims := CustomClaims{
		uid,
		jwt.StandardClaims{
			Id:        stringUid,
			Subject:   "miniapp",
			Audience:  "miniapp",
			ExpiresAt: time.Now().Unix() + (7 * 24 * 3600),
			Issuer:    "",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(this.JwtSecret))
	if err != nil {
		return "", errors.New(fmt.Sprintf(`create token err:%v`, err))
	}
	return tokenStr, err
}

func (this *gdrawUser) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		uid := c.Request.Header.Get("Uid")
		uidInt64, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": resp.RESP_CODE_NOAUTH_ERR, "msg": "未授权"})
			return
		}
		err = this.VerifyTokenFunc(uidInt64, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": resp.RESP_CODE_NOAUTH_ERR, "msg": err.Error()})
			return
		}
		c.Set("token", token)
		c.Set("uid", uidInt64)
		c.Next()
	}
}

func (this *gdrawUser) VerifyTokenFunc(uid int64, token string) (err error) {
	tokenValue, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(this.JwtSecret), nil
	})
	if err != nil {
		return
	}
	claims, ok := tokenValue.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("未授权")
		return
	}
	if claims.Uid != uid {
		err = errors.New("未授权")
		return
	}
	return
}

func (this *gdrawUser) Login(uid int, password string) (userinfo models.GdrawUser, loginSuccess bool, err error) {
	userinfo, _ = models.GdrawUserer.FindUserById(uid)
	if tools.FastFunc.MakeMd5(password) == userinfo.Password {
		loginSuccess = true
		return
	}
	return
}
