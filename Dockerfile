FROM golang:1.20-alpine AS builder
WORKDIR /go-chaos
COPY . . 
RUN go build -ldflags="-s -w" -o go-chaos .
FROM alpine:latest
COPY --from=builder /go-chaos/go-chaos /bin/go-chaos
ENTRYPOINT [ "/bin/go-chaos" ]

