exist := $(wildcard .envrc)
ifneq ($(strip $(exist)),)
  include .envrc
endif

.PHONY: run bin test audit tidy report install test/coverage

MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := go-sentry-tunnel

build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

run:
	@go run ./cmd/api/main.go

run/race:
	@go run -race ./cmd/api/main.go

bin:
	/tmp/bin/${BINARY_NAME}

test/coverage:
	@ENV=test go test -v ./... -coverprofile=coverage.out

test:
	@ENV=test go test -v -count=2 ./...

install:
	@go get -u ./...

report:
	@go tool cover -html=coverage.out -o coverage.html

audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

tidy:
	go fmt ./...
	go mod tidy -v
