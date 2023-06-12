FROM golang:alpine AS builder
RUN apk add git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . /go/src/wwwredirect

RUN cd /go/src/wwwredirect && go build -o www.redirect
FROM alpine
WORKDIR /api

COPY --from=builder /go/src/wwwredirect/www.redirect /api/

ENTRYPOINT ./www.redirect