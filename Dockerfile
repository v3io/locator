FROM golang:1.10 as builder

# copy sources
ADD . /go/src/github.com/v3io/locator

# build the locator
RuN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -o /go/bin/locatorctl src/github.com/v3io/locator/cmd/locatorctl/main.go

FROM debian:stretch-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

COPY --from=builder /go/bin/locatorctl /usr/local/bin/locatorctl

CMD [ "locatorctl" ]