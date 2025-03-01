FROM alpine:3.21.2 AS cert
RUN apk add --update --no-cache ca-certificates

FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /microsim

FROM scratch
COPY --from=cert /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /microsim /microsim
ENTRYPOINT ["/microsim"]
