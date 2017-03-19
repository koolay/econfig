FROM alpine:latest
WORKDIR /app
COPY econfig_linux_amd64 /app/econfig
COPY .econfig.toml /app/
ENTRYPOINT ["/app/econfig"]
