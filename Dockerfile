# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies for CGO (required for go-sqlite3)
# Using apk cache and trying different mirrors if needed
RUN apk add --no-cache gcc musl-dev sqlite-dev || \
    (echo "http://dl-4.alpinelinux.org/alpine/v3.22/main" > /etc/apk/repositories && \
     echo "http://dl-4.alpinelinux.org/alpine/v3.22/community" >> /etc/apk/repositories && \
     apk add --no-cache gcc musl-dev sqlite-dev)

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o aaxion ./cmd/main.go

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies or use fallback mirror
RUN apk --no-cache add ca-certificates sqlite-libs || \
    (echo "http://dl-4.alpinelinux.org/alpine/v3.18/main" > /etc/apk/repositories && \
     echo "http://dl-4.alpinelinux.org/alpine/v3.18/community" >> /etc/apk/repositories && \
     apk update && apk --no-cache add ca-certificates sqlite-libs)

# Create a non-root user
RUN addgroup -g 1000 aaxion && \
    adduser -D -u 1000 -G aaxion aaxion

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/aaxion .

# Create directories for uploads and data
RUN mkdir -p /data/uploads && \
    chown -R aaxion:aaxion /app /data

# Switch to non-root user
USER aaxion

# Expose port
EXPOSE 8080

# Set environment variables
ENV AAXION_DB_PATH=/data/.aaxion.db

# Run the application
CMD ["./aaxion"]
