# syntax = docker/dockerfile:experimental
FROM golang:1.12 as builder

WORKDIR /unbake

COPY *.go go.mod go.sum ./

ENV CGO_ENABLED 0
ENV GOOS linux

RUN --mount=type=cache,target=/root/.cache/go-build go test
RUN --mount=type=cache,target=/root/.cache/go-build go install

FROM alpine:3.10 as unbake
COPY --from=builder /go/bin/unbake /bin/unbake
CMD /bin/unbake
