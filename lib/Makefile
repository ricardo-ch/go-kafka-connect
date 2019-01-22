#!/usr/bin/env bash

.PHONY: install
install:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

.PHONY: test-integration
test-integration:
	make rundep
	go test -v `go list ./... | grep -v /vendor/` -tags=integration

.PHONY: rundep
rundep:
	docker-compose down
	docker-compose up -d
	@until $$(curl --output /dev/null --silent --head --fail http://localhost:8083); do \
		printf '.'; \
		sleep 2; \
	done
	@echo "up and running"

MOCKERY_PATH :=  $(shell  [ -z "$${GOBIN}" ] && echo $${GOPATH}/bin/mockery ||  echo $${GOBIN}/mockery; )

.PHONY: update-mocks
update-mocks:
	go get github.com/vektra/mockery/...
	mockery -inpkg -case "underscore" -recursive -all -note "NOTE: run 'make update-mocks' from this project top folder to update this file and generate new ones."

