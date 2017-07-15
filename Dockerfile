FROM ubuntu:16.04

ARG DAPPER_HOST_ARCH
ENV HOST_ARCH=${DAPPER_HOST_ARCH} ARCH=${DAPPER_HOST_ARCH}

RUN apt-get update && apt-get install -y pkg-config

RUN apt-get update && \
    apt-get install -y gcc make ca-certificates git wget curl vim less file && \
    rm -f /bin/sh && ln -s /bin/bash /bin/sh

ENV GOLANG_ARCH_amd64=amd64 GOLANG_ARCH_arm=armv6l GOLANG_ARCH=GOLANG_ARCH_${ARCH} \
    GOPATH=/go PATH=/go/bin:/usr/local/go/bin:${PATH} SHELL=/bin/bash

RUN wget -O - https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz | tar -xzf - -C /usr/local && \
    go get github.com/rancher/trash && go get github.com/golang/lint/golint
