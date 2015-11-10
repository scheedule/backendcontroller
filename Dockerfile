FROM ubuntu:14.04
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y golang git

# Set GOPATH
ENV GOPATH /go

# Grab Source
COPY . /go/src/github.com/scheedule/backendcontroller

WORKDIR /go/src/github.com/scheedule/backendcontroller

# Grab project dependencies and build
RUN go get ./... && go install

ENTRYPOINT ["/go/bin/backendcontroller"]
