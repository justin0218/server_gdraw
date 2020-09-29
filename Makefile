# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
#latest

run:
	export GOPROXY=http://goproxy.io
	$(GOCMD) run cmd/main.go

build:
	$(GOBUILD) -o gdraw cmd/main.go
