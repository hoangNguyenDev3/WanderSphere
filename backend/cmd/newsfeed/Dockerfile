FROM golang:1.20.5 AS builder

ENV GO111MODULE=on
WORKDIR /app

RUN mkdir /var/log/entry/
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/cache go mod download
RUN --mount=type=cache,target=/go/pkg/mod/cache CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o newsfeed.linux -mod=mod cmd/newsfeed/main.go
CMD ["/app/newsfeed.linux"]

FROM builder AS test_env
RUN cp /app/configs/config.yaml /app/config.yml
CMD ["/app/newsfeed.linux"]

FROM builder AS live_env
RUN cp /app/configs/live.yaml /app/config.yml
CMD ["/app/newsfeed.linux"]