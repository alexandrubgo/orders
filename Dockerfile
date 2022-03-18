ARG CONFIG

FROM golang:1.17-alpine AS builder
ARG NAME
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o $NAME /build

FROM alpine:3.13 AS certificates
RUN apk --no-cache add ca-certificates

# FROM scratch
FROM busybox:latest
ARG NAME
ARG CONFIG
ENV NAME=$NAME
ENV CONFIG=$CONFIG
WORKDIR /$NAME
ENV PATH=/orders/bin/:$PATH
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build ./
ENTRYPOINT ./$NAME service --config=$CONFIG