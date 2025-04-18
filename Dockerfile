FROM golang:1.23-alpine

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go app
RUN go build -o myapp .

# Expose port (sesuaikan dengan kebutuhan)
EXPOSE 8080

# Command to run the app
CMD ["./myapp"]
