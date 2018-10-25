# Stage 1: Build the binary
FROM golang:1.11.0-alpine as builder

# Install git
RUN apk update && apk add git

COPY . $GOPATH/src/github.com/gavincmartin/rotor-control-service/
WORKDIR $GOPATH/src/github.com/gavincmartin/rotor-control-service/

# get dependencies
RUN go get -d -v

# build the binary
RUN go build -o /go/bin/rotor-control-service

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# Stage 2: Build a smaller image from the alpine base
FROM alpine:3.8
MAINTAINER Gavin C. Martin

# Copy our static executable
COPY --from=builder /go/bin/rotor-control-service /app/rotor-control-service
COPY config.toml /app/

WORKDIR /app
ENTRYPOINT ["./rotor-control-service"]
