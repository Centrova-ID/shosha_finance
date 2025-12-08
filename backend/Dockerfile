# Build Stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server_binary ./cmd/cloud/main.go

# Run Stage
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server_binary .
EXPOSE 3000
CMD ["./server_binary"]
