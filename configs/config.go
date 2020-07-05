package configs

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var (
	Dft dfter
)

func init() {
	Dft = &dft{}
}

type dfter interface {
	Get() cfg
}

type dft struct {
	sync.Once
	dftCfg cfg
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
}

func (this *dft) Get() cfg {
	this.Do(func() {
		path := ""
		flag.StringVar(&path, "conf", "./configs/config.yml", "help")
		flag.Parse()
		bytes, err := ioutil.ReadFile(path)
		if nil != err {
			panic(err)
		}
		err = yaml.Unmarshal(bytes, &this.dftCfg)
		if nil != err {
			panic(err)
		}
	})
	return this.dftCfg
}
