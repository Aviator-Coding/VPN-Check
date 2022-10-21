############################
# STEP 1 build executable binary
############################
FROM golang:1.19-alpine AS builder 
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git 
#Change the Working Directory
WORKDIR $GOPATH/src/vpn-check-server/
#Clone Github Repo
COPY . .
# Fetch dependencies.
# And Compile Binary
RUN go get -d -v ./... &&\
    go install -v ./... &&\
    go build -o /go/bin/vpn-check-server

############################
# STEP 2 build a small image
############################
FROM alpine:edge
# Install SSL Certs
RUN apk add --no-cache ca-certificates
# Change Working dir
WORKDIR /vpn-check-server
# Copy our static executable.
COPY --from=builder /go/bin/vpn-check-server /usr/local/bin/
COPY --from=builder /go/src/vpn-check-server/ /vpn-check-server
EXPOSE 9000
#Run the Phantom Binary
ENTRYPOINT ["/usr/local/bin/vpn-check-server"]
