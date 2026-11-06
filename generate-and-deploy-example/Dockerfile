# Use an official Go image
FROM golang:1.22.3-alpine as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the web server
RUN go build -o server *.go

# Use a minimal image for the final container
FROM alpine:3.21.10@sha256:769c0a0cb2a64409697b54551072ee46ee1e36d4b369ddafd671b02a5d16e0e5
WORKDIR /app
COPY --from=builder /app/server ./

# Default entrypoint (can be overridden in docker-compose)
ENTRYPOINT ["./server"]