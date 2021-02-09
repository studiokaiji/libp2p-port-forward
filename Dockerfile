FROM golang:1.15.8

RUN apk update && apk add git

RUN mkdir /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app
