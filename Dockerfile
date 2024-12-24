# Start with a minimal image that supports Go binaries
FROM golang:1.21-alpine AS builder

# Install CA certificates
RUN apk update && apk add --no-cache ca-certificates tzdata

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files first to leverage Docker caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN go build -o server

# Use a minimal runtime image
FROM alpine:latest

# Install CA certificates for runtime too (important!)
RUN apk --no-cache add ca-certificates tzdata

# Copy the built binary from the builder stage
COPY --from=builder /app/server /server

# Expose the port Fly.io uses
EXPOSE 8080

# Command to run the binary
CMD ["./server"]

