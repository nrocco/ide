FROM alpine:3.6 AS builder
RUN apk --no-cache add bash automake autoconf build-base git ca-certificates openssl && \
    wget https://github.com/universal-ctags/ctags/archive/master.tar.gz && \
    tar xzf master.tar.gz && \
    cd ctags-master && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install

FROM alpine:3.6
RUN apk --no-cache update
COPY --from=builder /usr/local /usr/local
ENTRYPOINT ["ctags"]
