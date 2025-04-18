FROM golang:latest AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./ 

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/muj-submission-service ./submission_service/cmd/main.go

# Use a minimal alpine image for the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/muj-submission-service /app/muj-submission-service

# Copy CSV files
COPY Students_VIII.csv Students_VI.csv Student_IV.csv ./

# Expose the port
EXPOSE 8001

# Command to run
CMD ["/app/muj-submission-service"]