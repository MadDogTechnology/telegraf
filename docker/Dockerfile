
#
# docker rmi maddogtechnology-docker-develop.jfrog.io/telegraf
# docker build -t maddogtechnology-docker-develop.jfrog.io/telegraf:latest . -f docker/Dockerfile
# 

FROM tibetanmastiff/telegrafi-build as builder
COPY . /go/src/github.com/influxdata/telegraf/
RUN set -x \
    && cd $GOPATH/src/github.com/influxdata/telegraf \
    && make \
    && rm -rf /var/lib/apt/lists/*


FROM telegraf:1.9.3-alpine as app
COPY --from=builder /go/src/github.com/influxdata/telegraf/telegraf /usr/bin/telegraf

# Inherit from telegraf dockerfile https://github.com/influxdata/influxdata-docker/blob/master/telegraf/1.9/alpine/Dockerfile
# WORKDIR /
# CMD ["/entrypoint.sh "]

## Builder's notes:
#
#   FROM golang:alpine
#   RUN apk upgrade --update && apk add curl bash git && apk add ca-certificates alpine-sdk autoconf automake libtool dep && rm -f /var/cache/apk/*
#   CMD bash -c 'sleep 3153600000'
#