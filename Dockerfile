FROM golang:1.11 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/yurishkuro/microsim
WORKDIR /go/src/github.com/yurishkuro/microsim

COPY . .

RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -o /microsim

#FROM alpine:3.6
FROM scratch
COPY --from=builder /microsim /microsim
ENTRYPOINT ["/microsim"]
