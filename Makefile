export CGO_ENABLED := 0

TAGS := $(git describe --always)

build:
	@go build -ldflags="-s -w" .

image:
	@docker build -t zasdaym/echoer:$(TAGS) .
