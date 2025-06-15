# Build stage
FROM golang:1.23.6-alpine AS builder

# Enable Go modules
ENV GO111MODULE=on

# Set working directory inside container
WORKDIR /app

# Install git (required for Go to fetch private/public modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage - minimal image for running app
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Expose app port (adjust to your Gin port)
EXPOSE 8080

# Run the binary
CMD ["./main"]
