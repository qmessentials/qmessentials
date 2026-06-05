# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Install necessary build tools if needed (e.g., git for fetching dependencies)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY main.go ./

# Build the application
# CGO_ENABLED=0 produces a statically linked binary that works in minimal images
RUN CGO_ENABLED=0 GOOS=linux go build -o apigw main.go

# Stage 2: Runtime
# Using alpine for a small footprint while keeping basic shell utilities
FROM alpine:latest

# Install ca-certificates to ensure SSL/TLS connections work correctly
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/apigw .

# Expose the port the application listens on
EXPOSE 8080

# Set default environment variables
ENV PORT=8080
ENV GIN_MODE=release

# Note: API_SHARED_SECRET must be provided at runtime
# Note: INTAKE_URL and CONFIG_URL can also be provided at runtime

# Run the binary
CMD ["./apigw"]
