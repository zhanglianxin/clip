SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

CROSS_TARGETS := linux/amd64 darwin/amd64 windows/386 windows/amd64
OUT_PATH := out
BIN_FILE := clip


default: build

build:
	go build -ldflags "-X main.version=build_`date +%Y%m%d`" -o $(OUT_PATH)/$(BIN_FILE)

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=build_`date +%Y%m%d`" -o $(OUT_PATH)/$(BIN_FILE)_linux64

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=build_`date +%Y%m%d`" -o $(OUT_PATH)/$(BIN_FILE)_windows.exe

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -fr $(OUT_PATH)/*

rm-sha1:
	@rm -f $(OUT_PATH)/$(BIN_FILE)_*.sha1

gen-sha1: rm-sha1
	@$$(for f in $$(find $(OUT_PATH)/$(BIN_FILE)_* -type f); do shasum $$f > $$f.sha1; done)

copy-certs:
	cp -R certs $(OUT_PATH)/

zip: copy-certs
	zip $(BIN_FILE)_`date +%Y%m%d%H%M`.zip $(OUT_PATH)/

zipe: copy-certs
	zip $(BIN_FILE)_`date +%Y%m%d%H%M`.zip $(OUT_PATH)/ -e

tar: copy-certs
	tar -czvf $(BIN_FILE)_`date +%Y%m%d%H%M`.tar.gz $(BIN_FILE)/
