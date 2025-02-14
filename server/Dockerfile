# Use the official Golang image to create a build artifact.
FROM golang:1.23.4 as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code.
COPY . .

# Build the Go application.
RUN CGO_ENABLED=0 GOOS=linux go build -o go-api-template .

# Use a minimal Alpine image for the final stage.
FROM alpine:latest

# Set the working directory.
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/go-api-template .

# Expose the port the app runs on.
EXPOSE 8080

# Command to run the application.
CMD ["./go-api-template"]
