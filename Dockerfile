FROM golang:1.15.8

RUN apt-get update && apt-get install git

RUN mkdir /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app
