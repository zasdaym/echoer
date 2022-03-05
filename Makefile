export CGO_ENABLED := 0

TAGS := $(shell git describe --tags)

build:
	@go build -ldflags="-s -w -X main.version=${TAGS}" ./cmd/echoer/

compress:
	@upx echoer

image:
	@docker build -t ghcr.io/zasdaym/echoer:${TAGS} .
