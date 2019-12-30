FROM golang:1.13 as base

LABEL maintainer="Ivor Scott <ivor@devpie.io>"

WORKDIR /api

COPY . .

RUN go mod download

# Development
FROM base as dev

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src

RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /api

# port 4000 -> api port
# port 8888 -> debuggable api port
# port 2345 -> debugger port
EXPOSE 4000 8888 2345

# Production
FROM base as prod