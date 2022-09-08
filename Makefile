#!/bin/bash

CURDIR=$(shell pwd)
APP=`basename $(CURDIR)`
GO ?= go

all:
	#git pull
	go env -w GO111MODULE=on
	export GOPROXY=https://goproxy.io
	#mac-linux CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	$(GO) build -o bin/go-demo main.go

init:
	go get github.com/gin-gonic/gin
	go get github.com/BurntSushi/toml

clean:
	rm -rf bin/*
	rm -rf ./logs/*

.PHONY : clean
