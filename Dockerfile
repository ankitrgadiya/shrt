# Step 1: Build
FROM golang:1.21-alpine AS builder
RUN apk --update --no-cache add musl-dev gcc

WORKDIR /app
COPY . /app
RUN CC=/usr/bin/x86_64-alpine-linux-musl-gcc go build --ldflags '-linkmode external -extldflags "-static" -s -w' -o /shrt main.go

# Step 2: Final
FROM alpine:latest
COPY --from=builder /shrt /usr/local/bin/shrt
ENTRYPOINT ["/usr/local/bin/shrt"]
