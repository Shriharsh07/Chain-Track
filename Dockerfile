# Step 1: Use the official Go image to build the binary
FROM golang:1.22.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blockchain-api ./main.go

# Step 2: Use a minimal base image to run the app
FROM gcr.io/distroless/base-debian10

# Copy the compiled binary from the builder stage
COPY --from=builder /app/blockchain-api /blockchain-api

# Expose port (matching your Go server)
EXPOSE 8080

# Set the binary to run when the container starts
ENTRYPOINT ["/blockchain-api"]
