FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o zolaris-backend .

# Use a specific Alpine version for the final stage
FROM alpine:3.19

# Add necessary packages and create a non-root user
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/zolaris-backend /app/

# Set proper permissions
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/zolaris-backend"]
