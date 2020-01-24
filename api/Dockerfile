FROM golang:alpine3.10 as base

WORKDIR /api

FROM base as dev

COPY go.* ./

RUN go mod download

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get github.com/go-delve/delve/cmd/dlv \
    && go get github.com/githubnemo/CompileDaemon

# log go env (for debugging)
RUN go env

# port 4000 -> api port
# port 2345 -> debugger port

EXPOSE 4000 2345

FROM dev as test

COPY . .

RUN export CGO_ENABLED=0 && \
    go test -v ./...

FROM test as build-stage

RUN go build -ldflags "-s -w" -o main ./cmd/api

FROM base as prod

COPY --from=build-stage /api/main main

RUN addgroup gopher && adduser -D -G gopher gopher \
    && chown -R gopher:gopher /api && \
    chmod +x ./main

USER gopher

EXPOSE 4000
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:4000" ]

CMD ["./main"]
