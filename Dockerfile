FROM golang:1.23 AS builder

COPY . ./
LABEL authors="charlie.mills"

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o /app

FROM scratch
COPY --from=builder /app /whine

ENTRYPOINT ["/whine"]
