FROM golang:1.18 AS builder
WORKDIR /app
RUN apt-get -y update && apt-get -y install upx-ucl && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN make build compress

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/tmp/echoer .
CMD ["/app/echoer"]
