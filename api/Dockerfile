# 1. FROM sets the base image to use for subsequent instructions
# Use the golang alpine image as the base stage of a multi-stage routine
FROM golang:1.14-alpine as base

# 2. ENV CGO_ENABLED=0 diables CGO which is required for production builds
ENV CGO_ENABLED=0 

# 3. WORKDIR sets the working directory for any subsequent COPY, CMD, or RUN instructions
# Set the working directory to /api
WORKDIR /api

# 4. Extend the base stage and create a new stage named dev
FROM base as dev

# 5. COPY copies files or folders from source to the destination path in the image's filesystem
# Copy the go.mod and go.sum files to /api in the image's filesystem
COPY go.* ./

# 6. Install go module dependencies in the image's filesystem
RUN go mod download

# 7. ENV sets an environment variable
# Create GOPATH and PATH environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# 8. Print go environment for debugging purposes
RUN go env

# 9. Install development dependencies to debug and live reload api
RUN go get github.com/go-delve/delve/cmd/dlv \
    && go get github.com/githubnemo/CompileDaemon

# 10. Provide meta data about the ports the container must expose
# port 4000 -> api port
# port 2345 -> debugger port
EXPOSE 4000 2345

# 11. Extend the dev stage and create a new stage named build-stage
FROM dev as build-stage

# 12. Copy the remaining api code into /api in the image's filesystem
COPY . .

# 13. Build api
RUN go build -o main ./cmd/api

# 15. Install aquasecurity's trivy for robust image scanning
FROM aquasec/trivy:0.4.4 as trivy

# 16. Scan the alpine image before production use
RUN trivy alpine:3.11.5 && \
    echo "No image vulnerabilities" > result

# 17. Extend the base stage and create a new stage named prod
FROM alpine:3.11.5 as prod

# 18. Copy the build into the prod stage and leave everything else behind
COPY --from=trivy result secure
COPY --from=build-stage /api/main main

# 19. Provide meta data about the port the container must expose
EXPOSE 4000

# 20. Provide the default command for the production container
CMD ["./main"]
