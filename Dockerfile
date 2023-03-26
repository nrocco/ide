# syntax = docker/dockerfile:1-experimental
FROM --platform=${BUILDPLATFORM} golang:alpine AS godev
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        musl-dev \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
WORKDIR /src



FROM --platform=${BUILDPLATFORM} godev AS gobuilder
ENV CGO_ENABLED=0
COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
ARG BUILD_VERSION=master
ARG BUILD_COMMIT=unknown
ARG BUILD_DATE=now
ARG TARGETOS
ARG TARGETARCH
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build golint -set_exit_status ./...
RUN --mount=type=cache,target=/root/.cache/go-build go vet -v ./...
RUN --mount=type=cache,target=/root/.cache/go-build staticcheck ./...
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -short ./...
RUN mkdir -p dist
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -x -o dist \
        -ldflags "\
            -X github.com/nrocco/ide/cmd.version=${BUILD_VERSION} \
            -X github.com/nrocco/ide/cmd.commit=${BUILD_COMMIT} \
            -X github.com/nrocco/ide/cmd.date=${BUILD_DATE} \
            -s -w"
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -cover -short ./...



FROM scratch AS bin
COPY --from=gobuilder /src/dist/ /



FROM alpine:edge
RUN apk add --no-cache \
        ca-certificates \
        sqlite \
    && true
COPY --from=gobuilder /src/dist/ide /usr/bin/ide
ENTRYPOINT ["/usr/bin/ide"]
