# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY server/ ./server/

# Build the binary
WORKDIR /app/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o coap-server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server/coap-server .

# Run as non-root user
RUN addgroup -g 1000 coap && \
    adduser -D -u 1000 -G coap coap && \
    chown -R coap:coap /root

USER coap

# Expose CoAP port
EXPOSE 1373/udp

ENTRYPOINT ["./coap-server"]
