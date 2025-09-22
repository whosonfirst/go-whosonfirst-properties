FROM golang:1.25-alpine AS gotools

RUN mkdir /build

COPY . /build/go-whosonfirst-properties

RUN apk update && apk upgrade \
    && apk add git \
    #
    && ls -al /build/go-whosonfirst-properties \
    && cd /build/go-whosonfirst-properties \
    && go build -mod vendor -o /bin/index-properties cmd/index-properties/main.go \
    && cd - \
    #
    && git clone --depth https://github.com/aaronland/gocloud.git /build/gocloud \
    && cd /build/gocloud \
    && go build -mod vendor -o /bin/runtimevar cmd/runtimevar/main.go 
    
    
FROM alpine

RUN mkdir /usr/local/data

RUN apk update && apk upgrade \
    && apk add git bash

COPY --from=gotools /bin/index-properties /bin/index-properties
COPY --from=gotools /bin/runtimevar /bin/runtimevar

COPY docker-bin/index.sh /bin/index.sh