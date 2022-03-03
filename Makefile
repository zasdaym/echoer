export CGO_ENABLED := 0

COMMIT_SHA := $(shell git describe --tags)

build:
	@go build -ldflags="-s -w -X main.version=${COMMIT_SHA}" .

image:
	@docker build -t zasdaym/echoer:${COMMIT_SHA} .
