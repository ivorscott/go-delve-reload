# 1. FROM sets the base image to use for subsequent instructions
# Use the golang alpine image as the base stage of a multi-stage routine
FROM golang:1.14-alpine as base

# 2. WORKDIR sets the working directory for any subsequent COPY, CMD, or RUN instructions
# Set the working directory to /api
WORKDIR /api

# 3. Extend aquasecurity's trivy image and create a new stage named trivy
# Used for robust image scanning
FROM aquasec/trivy:0.4.4 as trivy 

# 4. RUN executes commands on top of the current image as a new layer and commits the results
# Scan the golang alpine image before production use
RUN trivy --debug --timeout 4m golang:1.14-alpine && \
    echo "No image vulnerabilities" > result

# 5. Extend the base stage and create a new stage named dev
FROM base as dev

# 6. COPY copies files or folders from source to the destination path in the image's filesystem
# Copy the go.mod and go.sum files to /api in the image's filesystem
COPY go.* ./

# 7. Install go module dependencies in the image's filesystem
RUN go mod download

# 8. ENV sets an environment variable
# Create GOPATH and PATH environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# 9. Print go environment for debugging purposes
RUN go env

# 10. Install development dependencies to debug and live reload api
RUN go get github.com/go-delve/delve/cmd/dlv \
    && go get github.com/githubnemo/CompileDaemon


# 11. Provide meta data about the ports the container must expose
# port 4000 -> api port
# port 2345 -> debugger port
EXPOSE 4000 2345

# 12. Extend the dev stage and create a new stage named test
FROM dev as test

# 13. Copy the remaining api code into /api in the image's filesystem
COPY . .

# 14. Disable CGO and run unit tests
RUN export CGO_ENABLED=0 && \
    go test -v ./...

# 15. Extend the test stage and create a new stage named build-stage
FROM test as build-stage

# 16. Build the api with "-ldflags" aka linker flags to reduce binary size
# -s = disable symbol table
# -w = disable DWARF generation
RUN GOOS=linux go build -ldflags "-s -w" -o main ./cmd/api

# 17. Extend the base stage and create a new stage named prod
FROM base as prod

# 18. Copy only the files we want from a few stages into the prod stage
COPY --from=trivy result secure
COPY --from=build-stage /api/main main

# 19. Create a new group and user, recursively change directory ownership, then allow the binary to be executed
RUN addgroup gopher && adduser -D -G gopher gopher \
    && chown -R gopher:gopher /api && \
    chmod +x ./main

# 20. Change to a non-root user
USER gopher

# 21. Provide meta data about the port the container must expose
EXPOSE 4000

# 22. Define how Docker should test the container to check that it is still working
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:4000" ]

# 23. Provide the default command for the production container
CMD ["./main"]