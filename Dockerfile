# Use the official Golang image to create a binary
FROM golang:1.20-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a smaller base image for the final stage
FROM alpine:latest
# Copy the binary from the builder stage to the final stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
