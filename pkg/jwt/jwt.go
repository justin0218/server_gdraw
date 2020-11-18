package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"server_gdraw/store"
	"strconv"
	"time"
)

var secret = new(store.Config).Get().Jwt.Secret

type CustomClaims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

func CreateToken(uid int64) (string, error) {
	stringUid := strconv.FormatInt(uid, 10)
	claims := CustomClaims{
		uid,
		jwt.StandardClaims{
			Id:        stringUid,
			Subject:   secret,
			Audience:  secret,
			ExpiresAt: time.Now().Unix() + (24 * 3600 * 7),
			Issuer:    "",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New(fmt.Sprintf(`create token err:%v`, err))
	}
	return tokenStr, err
}

func VerifyToken(tokenString string) (uid int64, err error) {
	tokenValue, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		err = errors.New(err.Error())
		return
	}
	claims, ok := tokenValue.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("token is invalid")
		return
	}
	uid = claims.Uid
	return
}

func GetUid(c *gin.Context) int {
	if val, ex := c.Get("uid"); ex {
		return int(val.(int64))
	}
	return 0
}
