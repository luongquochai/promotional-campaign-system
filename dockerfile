# Use official Golang image as a builder
FROM golang:1.18 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application
COPY . .

# Build the Go app
RUN go build -o server cmd/server/main.go

# Start a new stage from a smaller image
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Copy the Go binary from the builder
COPY --from=builder /app/server /server

# Expose the port your application will run on
EXPOSE 8080

# Command to run the application
CMD ["/server"]
