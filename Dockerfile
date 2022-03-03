FROM golang:1.17.6 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN make build

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/echoer .
CMD ["/app/echoer"]
