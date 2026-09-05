# Use an official Go image
FROM golang:1.27.1-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 as builder

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