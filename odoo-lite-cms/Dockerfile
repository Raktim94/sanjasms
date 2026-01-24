# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
# CGO_ENABLED=0 is preferred for static binaries, but sqlite might need it if not using the pure go version.
# glebarez/sqlite is pure Go, so CGO_ENABLED=0 is fine.
RUN CGO_ENABLED=0 GOOS=linux go build -o cms main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (check if needed for pure go)
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/cms .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Create directory for sqlite db
RUN mkdir -p /data

EXPOSE 8080

CMD ["./cms"]
