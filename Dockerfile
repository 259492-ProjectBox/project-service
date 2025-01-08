# Stage 1: Builder
# FROM golang:latest
FROM golang:alpine AS builder
# Install Poppler (includes pdftotext)
# RUN apt-get update && apt-get install -y poppler-utils
RUN apk add --no-cache poppler-utils
# Install air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Install swag for generating swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download


# Copy the entire project
COPY . .

# Generate Swagger documentation
RUN swag init --output ./docs

# Remove the `swag` binary to save space
# RUN rm -rf /go/bin/swag


# Stage 2: Runtime
FROM golang:alpine

# Install runtime dependencies
RUN apk add --no-cache poppler-utils

# Set working directory
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app /app
COPY --from=builder /go/bin/air /go/bin/air
# Set PATH for binaries
ENV PATH=$PATH:/go/bin

# Command to start the application
CMD ["air"]
