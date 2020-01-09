FROM golang:1.13.5 as base

WORKDIR /api

COPY go.* ./

RUN go mod download

FROM base as dev

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src

RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /api

# port 4000 -> api port
# port 2345 -> debugger port

EXPOSE 4000 2345

FROM dev as test

COPY . .

RUN go test -v ./...

FROM test as build-stage

RUN go build -o main ./cmd/api

FROM scratch as prod

COPY --from=build-stage /api/main main

CMD ["./main"]



