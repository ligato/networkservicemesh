FROM golang:1 as build
ENV PACKAGEPATH=github.com/ligato/networkservicemesh/

ENV http_proxy ${HTTP_PROXY}
ENV https_proxy ${HTTPS_PROXY}
RUN apt update
RUN apt install -y git bash unzip
RUN wget -e use_proxy=yes https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
RUN unzip -d /usr/local protoc-3.5.1-linux-x86_64
COPY [".","/go/src/${PACKAGEPATH}"]
WORKDIR /go/src/${PACKAGEPATH}/
RUN ./scripts/build.sh
