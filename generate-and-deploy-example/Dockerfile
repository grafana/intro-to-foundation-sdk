# Use an official Go image
FROM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb as builder

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