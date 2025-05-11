# Step 1: Use the official Go image to build the binary
FROM golang:1.23.9 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

COPY docs ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blockchain-api ./main.go

# Step 2: Use a minimal base image to run the app
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /app/blockchain-api /blockchain-api
COPY --from=builder /app/docs /docs

# Expose port (matching your Go server)
EXPOSE 8080

# Set the binary to run when the container starts
ENTRYPOINT ["/blockchain-api"]
