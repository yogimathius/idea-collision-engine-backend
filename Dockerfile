# Build stage for Go backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate cmd/migrate/main.go

# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# Copy package files
COPY package*.json pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build frontend
RUN pnpm run build

# Production stage
FROM alpine:3.20

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy backend binary and migration tool
COPY --from=backend-builder /app/main .
COPY --from=backend-builder /app/migrate .

# Copy frontend build
COPY --from=frontend-builder /frontend/dist ./static

# Copy migrations and documentation
COPY migrations ./migrations
COPY docs ./docs
COPY internal/handlers/swagger-ui ./internal/handlers/swagger-ui

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -D -s /bin/sh -u 1000 -G appgroup appuser

# Change ownership
RUN chown -R appuser:appgroup /root

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./main"]