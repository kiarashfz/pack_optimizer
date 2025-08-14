# Stage 1: The builder stage.
# We use a specific version for reproducibility and security.
FROM golang:1.25.0-bookworm AS builder

WORKDIR /app

# Copy only the go module files first to leverage Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the Go application as a static, non-CGO-enabled binary for portability.
# This results in a self-contained executable that doesn't need C libraries.
# We also use the `-ldflags` for a smaller binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/pack_optimizer ./cmd/server/main.go

# --- Stage 2: The final runtime image. ---
# Use a minimal Alpine Linux image for a tiny and secure final image.
FROM alpine:3.22.1

# Install ca-certificates and create a non-root user in a single layer.
RUN apk add --no-cache ca-certificates \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

USER appuser

# Set the working directory to a path accessible by the new user.
WORKDIR /app

# Copy the statically-built binary from the 'builder' stage.
COPY --from=builder --chown=appuser:appgroup /app/pack_optimizer .

# Copy the .env file for local development
COPY --from=builder --chown=appuser:appgroup /app/.env ./

# Expose the port on which the application listens.
EXPOSE 8080

# The command to run the application binary.
CMD ["./pack_optimizer"]
