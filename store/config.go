package store

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
}

type cfg struct {
	Runmode string `yaml:"runmode"`
	Grpc    struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	Http struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	RocketMQ struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKey       string `yaml:"accessKey"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		Topic           string `yaml:"topic"`
		InstanceId      string `yaml:"instanceId"`
		GroupId         string `yaml:"groupId"`
	} `yaml:"rocketMQ"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Pass string `yaml:"password"`
	} `yaml:"redis"`
	Wechat struct {
		Appid              string `yaml:"appid"`
		Secret             string `yaml:"secret"`
		MchApiKey          string `yaml:"mchApiKey"`
		MchID              string `yaml:"mchId"`
		BillNoticeTemplate string `yaml:"billNoticeTemplate"`
		BillNoticeUrl      string `yaml:"billNoticeUrl"`
	} `yaml:"wechat"`
	Mysql struct {
		Master struct {
			Addr         string `yaml:"addr"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
		} `yaml:"master"`
		AppRead struct {
			Addr         string `yaml:"addr"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
		} `yaml:"appRead"`
		AdminRead struct {
			Addr         string `yaml:"addr"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
		} `yaml:"adminRead"`
	} `yaml:"mysql"`
	Jwt struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	Miniapp struct {
		Appid  string `yaml:"appid"`
		Secret string `yaml:"secret"`
	} `yaml:"miniapp"`
}

func (s *Config) Get() (conf cfg) {
	configOnce.Do(func() {
		path := ""
		flag.StringVar(&path, "conf", "./configs/config.yml", "help")
		flag.Parse()
		bytes, err := ioutil.ReadFile(path)
		if nil != err {
			panic(err)
		}
		err = yaml.Unmarshal(bytes, &conf)
		if nil != err {
			panic(err)
		}
		config = conf
	})
	conf = config
	return
}
