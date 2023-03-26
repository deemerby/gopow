# build the app binary
FROM golang:1.19 AS builder
ENV GO111MODULE=on
WORKDIR /go/src/github.com/deemerby/gopow
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -mod=vendor -o bin/app ./cmd/app/

# copy the app binary from builder stage; run the app binary
FROM alpine:latest AS runner
WORKDIR /bin

# Go programs require libc
RUN mkdir -p /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder /go/src/github.com/deemerby/gopow/bin/app .

ENTRYPOINT ["app"]
