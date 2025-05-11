FROM golang:1.20.5 AS builder

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o newsfeed.linux -mod=mod cmd/newsfeed/main.go

FROM alpine:latest AS test_env
WORKDIR /app
COPY --from=builder /app/newsfeed.linux .
COPY --from=builder /app/configs /app/configs
CMD ["/app/newsfeed.linux", "-conf", "configs/files/test.yml"]

FROM alpine:latest AS live_env
WORKDIR /app
COPY --from=builder /app/newsfeed.linux .
COPY --from=builder /app/configs /app/configs
CMD ["/app/newsfeed.linux", "-conf", "configs/files/live.yml"]