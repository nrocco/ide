# syntax = docker/dockerfile:1-experimental
FROM --platform=${BUILDPLATFORM} golang:alpine AS gobase
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        musl-dev \
    && true
RUN env GO111MODULE=on go get -u \
        golang.org/x/lint/golint \
        golang.org/x/tools/cmd/goimports \
    && true
WORKDIR /src



FROM --platform=${BUILDPLATFORM} gobase AS gobuilder
ENV CGO_ENABLED=0
COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
ARG BUILD_VERSION=master
ARG BUILD_COMMIT=unknown
ARG BUILD_DATE=now
ARG TARGETOS
ARG TARGETARCH
COPY . .
# RUN protoc -I server/ server/server.proto --go_out=plugins=grpc:server
RUN --mount=type=cache,target=/root/.cache/go-build golint -set_exit_status ./...
RUN --mount=type=cache,target=/root/.cache/go-build go vet -v ./...
RUN mkdir -p dist
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -x -o dist \
        -ldflags "\
            -X github.com/nrocco/ide/cmd.version=${BUILD_VERSION} \
            -X github.com/nrocco/ide/cmd.commit=${BUILD_COMMIT} \
            -X github.com/nrocco/ide/cmd.date=${BUILD_DATE} \
            -s -w"
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -short ./...



FROM scratch AS bin
COPY --from=gobuilder /src/dist/ /



FROM alpine:edge
RUN apk add --no-cache \
        ca-certificates \
        sqlite \
    && true
COPY --from=gobuilder /src/dist/ide /usr/bin/ide
EXPOSE 3000
WORKDIR /var/lib/ide
VOLUME /var/lib/ide
CMD ["ide", "server"]
