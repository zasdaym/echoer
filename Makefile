export CGO_ENABLED := 0

VERSION := $(shell git describe --tags)

build:
	@go build -ldflags="-s -w -X main.version=${VERSION}" ./cmd/echoer/

compress:
	@upx echoer

image:
	@docker build -t ghcr.io/zasdaym/echoer:${VERSION} .
