package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

var FastFunc fastFuncer

type fastFunc struct {
}

type fastFuncer interface {
	MakeMd5(str string) (md5str string)
	DataToString(d interface{}) (datastr string)
	DataToBytes(d interface{}) (databyte []byte)
}

func init() {
	FastFunc = &fastFunc{}
}

func (this *fastFunc) MakeMd5(str string) (md5str string) {
	h := md5.New()
	h.Write([]byte(str))
	md5str = hex.EncodeToString(h.Sum(nil))
	return
}

func (this *fastFunc) DataToString(d interface{}) (datastr string) {
	da, _ := json.Marshal(d)
	datastr = string(da)
	return
}

func (this *fastFunc) DataToBytes(d interface{}) (databyte []byte) {
	da, _ := json.Marshal(d)
	databyte = da
	return
}
