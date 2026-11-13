# Use an official Go image
FROM golang:1.27.4-alpine@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf as builder

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