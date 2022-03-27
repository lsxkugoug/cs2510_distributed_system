# syntax=docker/dockerfile:1
FROM golang:1.17.6

ADD . /p2

# WORKDIR /p2/code/
# RUN go install

CMD ["/bin/sh", "-c","cd /p2/code/src; go test -v -run Test1;"]
