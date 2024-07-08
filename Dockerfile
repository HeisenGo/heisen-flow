# Use the official Golang image to build the Go application
FROM golang:1.22.2 AS builder

# Set the Current Working Directory inside the container
WORKDIR /

# Copy go mod and sum files
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
#COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/api

# Use a smaller image for the final stage
FROM scratch

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the pre-built binary file from the builder image
COPY --from=builder /app /

# Copy the config file
COPY --from=builder /config.yaml /

# Expose port 2244 to the outside world
EXPOSE 2244

# Command to run the executable
CMD ["./app", "--config", "config.yaml"]

