FROM debian:stretch-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

COPY locatorctl /usr/local/bin/locatorctl

CMD [ "locatorctl" ]