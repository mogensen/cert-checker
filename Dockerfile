FROM golang:latest  AS build-env

# Dependencies
WORKDIR /build
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /app/cert-checker

# Build runtime container
FROM scratch
LABEL description="Certificate monitoring utility for watching tls certificates and reporting the result as metrics."
WORKDIR /app
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env --chown=1000:1000 /app /app

USER 1000:1000

CMD ["/app/cert-checker", "-f", "/data/config.yaml", "--apply"]
