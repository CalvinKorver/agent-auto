# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files from backend directory
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
