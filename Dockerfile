# Stage 1: Build the binary
FROM golang:1.11.0-alpine as builder

# Install git
RUN apk update && apk add git && apk add tzdata

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

# Install CA Certificates and copy timezone info from the builder
RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app
ENTRYPOINT ["./rotor-control-service"]
