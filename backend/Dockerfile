FROM golang:1.22.6-alpine3.20

WORKDIR /go/src/tweet-app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /go/bin/tweet-app ./...