FROM golang:alpine as build

RUN apk add --no-cache ca-certificates build-base

WORKDIR /build

ADD . .

# CGO enabled for Herumi BLS
RUN CGO_ENABLED=1 GOOS=linux \
    go build -ldflags '-extldflags "-static"' -o app

FROM scratch

# For future p2p/api commands
COPY --from=build /etc/ssl/certs/ca-certificates.crt \
     /etc/ssl/certs/ca-certificates.crt

COPY --from=build /build/app /app

ENTRYPOINT ["/app"]
