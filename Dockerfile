# Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN go build -o webhook .

FROM alpine:3.23

RUN apk add --no-cache ca-certificates
COPY --from=builder /workspace/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook"]