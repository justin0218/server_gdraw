package thirdParty

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

var TaoBao taobaoer

type taobaoer interface {
	GetIpDetail(ip string) (detail ipdetail)
}

type taobao struct {
	ipdetailUrl string
}

func init() {
	TaoBao = &taobao{
		ipdetailUrl: "http://ip.taobao.com/service/getIpInfo.php",
	}
}

type ipdetail struct {
	Data struct {
		Country string `json:"country"`
		Region  string `json:"region"`
		City    string `json:"city"`
		Isp     string `json:"isp"`
	} `json:"data"`
}

func (this *taobao) GetIpDetail(ip string) (detail ipdetail) {
	rurl := fmt.Sprintf("%s?ip=%s", this.ipdetailUrl, ip)
	gorequest.New().Get(rurl).Timeout(time.Second * 30).EndStruct(&detail)
	return
}
