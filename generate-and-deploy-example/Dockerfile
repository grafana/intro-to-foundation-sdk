# Use an official Go image
FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the web server
RUN go build -o server *.go

# Use a minimal image for the final container
FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder /app/server ./

# Default entrypoint (can be overridden in docker-compose)
ENTRYPOINT ["./server"]