FROM golang:latest AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/muj-portal-service ./portal_service/cmd/main.go

# Use a minimal alpine image for the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/muj-portal-service /app/muj-portal-service

# Expose the port
EXPOSE 8002

# Command to run
CMD ["/app/muj-portal-service"]