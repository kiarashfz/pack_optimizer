# Stage 1: The builder stage.
# We use a specific version for reproducibility and security.
FROM golang:1.23.4-bullseye AS builder

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
FROM alpine:3.18.4

# Install ca-certificates for secure HTTPS connections.
RUN apk add --no-cache ca-certificates

# Create a non-root user and group to run the application securely.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the working directory to a path accessible by the new user.
WORKDIR /app

# Copy the statically-built binary from the 'builder' stage.
COPY --from=builder --chown=appuser:appgroup /app/pack_optimizer .

# Copy the necessary templates and database migrations from the builder stage.
# We need migrations at runtime for the 'go-migrate' tool to work.
COPY --from=builder --chown=appuser:appgroup /app/templates ./templates
COPY --from=builder --chown=appuser:appgroup /app/db/migrations ./db/migrations

# Expose the port on which the application listens.
EXPOSE 8080

# The command to run the application binary.
CMD ["./pack_optimizer"]
