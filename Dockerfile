FROM golang:latest as builder

ADD . /go/src/github.com/MofeLee/reinit-cluster
WORKDIR /go/src/github.com/MofeLee/reinit-cluster
RUN go get -u github.com/golang/dep/cmd/dep \
  && dep ensure \
  && go install github.com/MofeLee/reinit-cluster

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/reinit-cluster .
