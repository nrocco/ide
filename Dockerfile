# syntax = docker/dockerfile:1
FROM --platform=${BUILDPLATFORM} tonistiigi/xx:latest AS xx

FROM --platform=${BUILDPLATFORM} golang:alpine AS godev
COPY --from=xx / /
RUN apk add --no-cache \
        ca-certificates \
        clang \
        ctags \
        file \
        findutils \
        gcc \
        musl-dev \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/cmd/deadcode@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install golang.org/x/vuln/cmd/govulncheck@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
WORKDIR /src
ENV CGO_ENABLED=0

FROM godev AS govendorer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod go mod download

FROM govendorer AS gobuilder
ARG BUILD_VERSION=master
ARG BUILD_COMMIT=unknown
ARG BUILD_DATE=now
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod golint -set_exit_status ./...
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod go vet -v ./...
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod staticcheck ./...
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod govulncheck ./...
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod deadcode .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod go test -v -cover -short ./...
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod \
    mkdir -p /out && \
    xx-go build -trimpath -o /out -ldflags "\
        -X github.com/nrocco/ide/cmd.version=${BUILD_VERSION} \
        -X github.com/nrocco/ide/cmd.commit=${BUILD_COMMIT} \
        -X github.com/nrocco/ide/cmd.date=${BUILD_DATE} \
        -s -w" && \
    xx-verify --static /out/*

FROM scratch AS bin
COPY --from=gobuilder /out/ /
