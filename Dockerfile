FROM golang:1.12 as builder

WORKDIR /unbake

COPY *.go go.mod go.sum ./

ENV CGO_ENABLED 0
ENV GOOS linux

RUN go test
RUN go install

FROM alpine:3.10 as unbake
COPY --from=builder /go/bin/unbake /bin/unbake
CMD /bin/unbake
