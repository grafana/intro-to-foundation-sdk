# Use an official Go image
FROM golang:1.27.4-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 as builder

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