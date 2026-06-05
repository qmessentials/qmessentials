# Build stage
FROM golang:1.26.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# Use CGO_ENABLED=0 to build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates in case the app makes HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Expose the default port
EXPOSE 8082

# Define environment variables with placeholders or defaults if applicable
# Note: These should be provided at runtime
ENV PORT=8082
# ENV API_SHARED_SECRET=
# ENV DB_HOST=
# ENV DB_PORT=
# ENV DB_NAME=
# ENV DB_USER=
# ENV DB_PASSWORD=

# Command to run the executable
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

CMD ["./main"]
