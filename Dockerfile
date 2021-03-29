#####
FROM golang:1.16.2 as base_api
WORKDIR /go/app

ENV GO111MODULE=off
RUN set -eux \
    && apt-get -y update \
    && apt-get -y install git sqlite3 curl \
    && go get github.com/go-delve/delve/cmd/dlv  \
    && go build -o /go/bin/dlv github.com/go-delve/delve/cmd/dlv \
    && go get github.com/cosmtrek/air \
    && go build -o /go/bin/air github.com/cosmtrek/air
ENV GO111MODULE on

RUN sqlite3 /tmp/scapo.db


#####
FROM base_api as build_api
WORKDIR /go/app
COPY . /go/app/

RUN sqlite3 /tmp/scapo.db < sql/setup.sql

CMD air -c .air.toml


