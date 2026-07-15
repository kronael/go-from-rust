FROM golang:1.26.0-alpine3.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -tags web -o /out/go-from-rust-web .

FROM alpine:3.23.3
RUN apk add --no-cache ca-certificates \
    && addgroup -S app && adduser -S -G app app
COPY --from=build /out/go-from-rust-web /usr/local/bin/go-from-rust-web
USER app
ENV ADDR=:8080
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O- http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/usr/local/bin/go-from-rust-web"]
