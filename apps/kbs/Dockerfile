# Builder image go
FROM golang:1.21.0 AS builder

ARG appVersion
ARG commitHash

ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash

# Build kbs binary with Go
ENV GOPATH /opt/go

RUN mkdir -p /kbs
WORKDIR /kbs
COPY . /kbs
RUN go mod tidy
RUN make build-linux

# Runnable image
FROM alpine
ARG appVersion
ARG commitHash
ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash
COPY --from=builder /kbs/bin/kbs-amd64-linux /bin/kbs-service
RUN ls /bin/kbs-service
WORKDIR /bin
ENTRYPOINT [ "./kbs-service" ]