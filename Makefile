# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
#latest
export GOPROXY=http://goproxy.cn
run:
	$(GOCMD) run cmd/main.go

build:
	$(GOBUILD) -o gdraw cmd/main.go
