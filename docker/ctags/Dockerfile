FROM alpine:edge AS builder
RUN apk --no-cache add bash automake autoconf build-base git ca-certificates openssl && \
    wget https://github.com/universal-ctags/ctags/archive/master.tar.gz && \
    tar xzf master.tar.gz && \
    cd ctags-master && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install

FROM alpine:edge
RUN apk --no-cache update
COPY --from=builder /usr/local /usr/local
ENTRYPOINT ["ctags"]
