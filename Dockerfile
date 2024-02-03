FROM debian:bullseye-slim

LABEL maintainer="merdw"

# Set the Go version
ARG GO_VERSION=1.20.13
ARG LOAD_LANG=tur

# Install necessary dependencies
RUN apt update \
    && apt install -y \
        ca-certificates \
        libtesseract-dev=4.1.1-2.1 \
        tesseract-ocr=4.1.1-2.1 \
        wget \
        build-essential \
    && if [ -n "${LOAD_LANG}" ]; then apt-get install -y tesseract-ocr-${LOAD_LANG}; fi \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=${HOME}/go
ENV PATH=${PATH}:${GOPATH}/bin

## Add the current directory to the Go workspace
ADD . $GOPATH/src/github.com/merdw/fiberocr
WORKDIR $GOPATH/src/github.com/merdw/fiberocr

ADD . .

RUN go build  -o $GOPATH/src/github.com/merdw/fiberocr/temp
