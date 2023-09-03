# syntax=docker/dockerfile:1.2

FROM golang:1.21rc4-alpine3.18

RUN go version
ENV GOPATH=/

COPY ./ /binance-converter-telegram-bot

WORKDIR /binance-converter-telegram-bot

# build go app
RUN go mod download

RUN go mod tidy

RUN go build -o binance-converter-telegram-bot ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates


WORKDIR /root/

COPY --from=0 /binance-converter-telegram-bot/binance-converter-telegram-bot .
