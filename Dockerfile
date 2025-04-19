FROM golang:latest AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy common packages and directories needed by both services
COPY pkg ./pkg
COPY migrations ./migrations

# Build portal service
COPY portal_service ./portal_service
RUN CGO_ENABLED=0 GOOS=linux go build -o muj-portal-service ./portal_service/cmd/main.go

# Build submission service
COPY submission_service ./submission_service
RUN CGO_ENABLED=0 GOOS=linux go build -o muj-submission-service ./submission_service/cmd/main.go

# Create final minimal image
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/muj-portal-service ./muj-portal-service
COPY --from=builder /app/muj-submission-service ./muj-submission-service
COPY --from=builder /app/migrations ./migrations

# Copy CSV files
COPY Students_VIII.csv Students_VI.csv Student_IV.csv ./

EXPOSE 8001 8002

# The CMD will be specified in docker-compose.yml
