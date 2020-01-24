# 1. FROM sets the baseImage to use for subsequent instructions
# Use the Golang image as the base stage of a multi-stage routine
FROM golang:1.13.5 as base

# 2. WORKDIR sets the working directory for any subsequent COPY, CMD, or RUN instructions
# Change directory
WORKDIR /api

# 3. Extend the base stage to create a new stage named dev
FROM base as dev

# 4. COPY copies files or folders from source to the destination path in the image's filesystem
# Copy the go.mod and go.sum files to /api in the image's filesystem
COPY go.* ./

# 5. RUN executes commands on top of the current image as a new layer and commit the results
# Install go module dependencies in the image's filesystem
RUN go mod download

# 6. ENV sets an environment variable
# Create GOPATH and PATH Environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# 7. Log go environment (for debugging)
RUN go env

# 8. Install development dependencies to debug and live reload api
RUN go get github.com/go-delve/delve/cmd/dlv \
    && go get github.com/githubnemo/CompileDaemon


# 9. Provide meta data about the ports the container must expose
# port 4000 -> api port
# port 2345 -> debugger port
EXPOSE 4000 2345

# 10. Extend the dev stage to create a new stage named test
FROM dev as test

# 11. Copy the remaining api code into /api in the image's filesystem
COPY . .

# 12. Disable CGO and run unit tests
RUN export CGO_ENABLED=0 && \
    go test -v ./...

# 13. Extend the test stage to create a new stage named build-stage
FROM test as build-stage

# 14. Build the api with "-ldflags" aka linker flags to reduce binary size
# -s = disable symbol table
# -w = disable DWARF generation
RUN go build -ldflags "-s -w" -o main ./cmd/api

# 15. Extend the base stage to create a new stage named prod
FROM base as prod

# 16. Copy only the files we want from the build stage into the prod stage
COPY --from=build-stage /api/main main

# 17. Create a new group and user, recursively change directory ownership, then allow the binary to be executed
RUN addgroup gopher && adduser -D -G gopher gopher \
    && chown -R gopher:gopher /api && \
    chmod +x ./main

# 18. Change to a non-root user
USER gopher

# 19. Provide meta data about the port the container must expose
EXPOSE 4000

# 20. Define how Docker should test the container to check that it is still working
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:4000" ]

# 21. Provide the default command for the production container
CMD ["./main"]