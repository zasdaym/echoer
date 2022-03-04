export CGO_ENABLED := 0

TAGS := $(shell git describe --tags)

build:
	@go build -ldflags="-s -w -X main.version=${TAGS}" .

image:
	@docker build -t zasdaym/echoer:${TAGS} .
