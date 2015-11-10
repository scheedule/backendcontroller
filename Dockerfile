FROM ubuntu:14.04
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y golang git

# Set GOPATH
ENV GOPATH /go

# Grab Source
COPY . /go/src/github.com/scheedule/backend_controller

WORKDIR /go/src/github.com/scheedule/backend_controller

# Grab project dependencies and build
RUN go get ./... && go install

ENTRYPOINT ["/go/bin/backend_controller"]
