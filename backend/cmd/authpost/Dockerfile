FROM golang:1.20.5 AS builder

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authpost.linux -mod=mod cmd/authpost/main.go

FROM alpine:latest AS test_env
WORKDIR /app
COPY --from=builder /app/authpost.linux .
COPY --from=builder /app/configs /app/configs
CMD ["/app/authpost.linux"]

FROM alpine:latest AS live_env
WORKDIR /app
COPY --from=builder /app/authpost.linux .
COPY --from=builder /app/configs /app/configs
CMD ["/app/authpost.linux"]