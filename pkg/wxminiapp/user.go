package wxminiapp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
)

type User struct {
	Openid    string `json:"openid"`
	Nickname  string `json:"nickName"`
	Avatar    string `json:"avatarUrl"`
	UnionId   string `json:"unionId"`
	Watermark struct {
		APPID string `appid`
	} `json:"watermark"`
}

//解密用户信息
func DecryptUserInfo(sessionKey, encryptedData, iv string) (userinfo User, err error) {
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
	aesPlantText = PKCS7UnPadding(aesPlantText)
	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	aesPlantText = []byte(re.ReplaceAllString(string(aesPlantText), "$1"))
	err = json.Unmarshal(aesPlantText, &userinfo)
	if err != nil {
		err = errors.New("DecodeJsonError" + err.Error())
		return
	}
	if userinfo.Watermark.APPID != appid {
		err = errors.New("appID is not match")
		return
	}
	return
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}
