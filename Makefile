# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
#latest
all: test build
update:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -o server_gdraw cmd/main.go
	mkdir cmd/linux
	mv ./server_gdraw cmd/linux
	scp cmd/linux/server_gdraw root@140.143.188.219:/www/bin/gdraw
	rm -rf cmd/linux
test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

download:
	export GOPROXY=https://goproxy.io
	$(GOCMD) mod download

run:
	export GOPROXY=http://goproxy.io
	$(GOCMD) run cmd/main.go

stop:
	pkill -f target/logic
	pkill -f target/job
	pkill -f target/comet
