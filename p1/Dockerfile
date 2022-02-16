# syntax=docker/dockerfile:1
FROM golang:1.17.6

ADD . /p1

# WORKDIR /p1/code/
# RUN go install

CMD ["/bin/sh", "-c","go version;cd /p1/code/src; go test -v -run Test1;go test -v -run Test2;go test -v -run Test3; go run printLog.go "]