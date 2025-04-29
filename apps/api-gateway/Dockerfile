# Use Go 1.21+ for the build stage
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o api-gateway

# Use Alpine-based final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/api-gateway .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./api-gateway"]
