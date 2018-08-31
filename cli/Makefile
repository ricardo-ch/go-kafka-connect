#!/usr/bin/env bash

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./kccli -a -ldflags '-s' -installsuffix cgo ./main.go

.PHONY: build-all
build-all:
	GOOS=linux CGO_ENABLED=0 go build -o ./bin/kccli -a -ldflags '-s' -installsuffix cgo ./main.go
	GOOS=darwin CGO_ENABLED=0 go build -o ./bin/kccli_mac -a -ldflags '-s' -installsuffix cgo ./main.go
	GOOS=windows CGO_ENABLED=0 go build -o ./bin/kccli.exe -a -ldflags '-s' -installsuffix cgo ./main.go

