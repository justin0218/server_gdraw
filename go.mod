module server_gdraw

go 1.12

require (
	github.com/astaxie/beego v1.12.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/goproxy v0.0.0-20200425205933-4ba4e08a9b7f // indirect
	github.com/gin-gonic/gin v1.6.2
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gogo/protobuf v1.1.1
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/itsjamie/gin-cors v0.0.0-20160420130702-97b4a9da7933
	github.com/jinzhu/gorm v1.9.12
	github.com/micro/protobuf v0.0.0-20180321161605-ebd3be6d4fdb
	github.com/mojocn/base64Captcha v1.3.1
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.8.0
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/yaml.v2 v2.2.8
	moul.io/http2curl v1.0.0 // indirect
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200423211502-4bdfaf469ed5
	golang.org/x/image => github.com/golang/image v0.0.0-20200430140353-33d19683fad8
	golang.org/x/mod => github.com/golang/mod v0.2.0
	golang.org/x/net => github.com/golang/net v0.0.0-20200421231249-e086a090c8fd
	golang.org/x/sync => github.com/golang/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/sys => github.com/golang/sys v0.0.0-20200420163511-1957bb5e6d1f
	golang.org/x/tools => github.com/golang/tools v0.0.0-20200425043458-8463f397d07c
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/appengine => github.com/golang/appengine latest
    golang.org/x/text => github.com/golang/text latest
)
