FROM golang:1.27.1-alpine3.23 AS build
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local
COPY go.mod web.go README.md [0-9][0-9]_*.go ./
COPY internal/webapp ./internal/webapp
COPY web ./web
RUN go build -trimpath -ldflags="-s -w" -tags web -o /out/go-from-rust-web .

FROM alpine:3.23.5
RUN apk add --no-cache ca-certificates \
    && addgroup -S app && adduser -S -G app app
COPY --from=build /out/go-from-rust-web /usr/local/bin/go-from-rust-web
USER app
# ADDR is the listen address (host networking binds it directly); the
# healthcheck reads the port straight from ADDR. The tour serves under a
# relative <base href>, so it needs no prefix config behind a proxy.
ENV ADDR=:49169
EXPOSE 49169
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null "http://127.0.0.1:${ADDR##*:}/health" || exit 1
ENTRYPOINT ["/usr/local/bin/go-from-rust-web"]
