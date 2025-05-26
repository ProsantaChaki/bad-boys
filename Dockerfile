# Build stage
FROM golang:1.23.4-alpine

WORKDIR /app

# Install git, bash and CompileDaemon
RUN apk add --no-cache git bash && \
    go install github.com/githubnemo/CompileDaemon@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 3003

# Run the application with CompileDaemon
CMD ["CompileDaemon", "--build=go build -o app ./cmd/server", "--command=./app"] 