
#
# docker rmi maddogtechnology-docker-develop.jfrog.io/telegraf
# docker build -t maddogtechnology-docker-develop.jfrog.io/telegraf:latest . -f docker/Dockerfile
# 

FROM tibetanmastiff/telegrafi-build as builder
LABEL telegraf=builder
COPY . /go/src/github.com/MadDogTechnology/telegraf/
RUN set -x \
    && cd $GOPATH/src/github.com/MadDogTechnology/telegraf \
    && make \
    && rm -rf /var/lib/apt/lists/*


FROM telegraf:1.9.3-alpine as app
COPY --from=builder /go/src/github.com/MadDogTechnology/telegraf/telegraf /usr/bin/telegraf
COPY --from=builder /go/src/github.com/MadDogTechnology/telegraf/docker/entrypoint.sh /go/src/github.com/MadDogTechnology/telegraf/docker/telegraf-template.conf /
 
RUN set -x \
    && apk --update --no-cache --virtual add \
       rsyslog \
       bash \
       curl \
       wget \
       gettext \
       jq \
    && rm -rf /var/cache/apk/*

# Inherit from telegraf dockerfile https://github.com/influxdata/influxdata-docker/blob/master/telegraf/1.9/alpine/Dockerfile
# WORKDIR /
# CMD ["/entrypoint.sh "]

## Builder's notes:
#
#   FROM golang:alpine
#   RUN apk upgrade --update && apk add curl bash git && apk add ca-certificates alpine-sdk autoconf automake libtool dep && rm -f /var/cache/apk/*
#   CMD bash -c 'sleep 3153600000'
#