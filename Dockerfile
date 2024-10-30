FROM golang:1.23 AS builder

COPY . ./
LABEL authors="charlie.mills"

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o /app

FROM builder AS test
RUN go test -v ./...

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/charlierm/whine/
COPY --from=builder /app /whine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/whine"]
