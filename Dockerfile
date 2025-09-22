FROM golang:1.25-alpine AS gotools

RUN mkdir /build

COPY . /build/go-whosonfirst-properties

RUN apk update && apk upgrade \
    && apk add git

RUN cd /build/go-whosonfirst-properties \
    && go build -mod vendor -ldflags="-s -w" -o /bin/index-properties cmd/index-properties/main.go

RUN cd /build \
    && git clone --depth 1 https://github.com/aaronland/gocloud.git \
    && cd gocloud \
    && go build -mod vendor -ldflags="-s -w" -o /bin/runtimevar cmd/runtimevar/main.go 
    
    
FROM alpine

RUN mkdir /usr/local/data

RUN apk update && apk upgrade \
    && apk add git bash

COPY --from=gotools /bin/index-properties /bin/index-properties
COPY --from=gotools /bin/runtimevar /bin/runtimevar

COPY docker-bin/index.sh /bin/index.sh