FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /microsim

#FROM alpine:3.6
FROM scratch
COPY --from=builder /microsim /microsim
ENTRYPOINT ["/microsim"]
