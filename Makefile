export CGO_ENABLED := 0

VERSION := $(shell git describe --tags)

build:
	@go build -ldflags="-s -w -X main.version=${VERSION}" -o ./tmp/echoer ./cmd/echoer/

compress: build
	@upx ./tmp/echoer
