# FROM golang:latest
FROM golang:alpine
# Install Poppler (includes pdftotext)
# RUN apt-get update && apt-get install -y poppler-utils
RUN apk add --no-cache poppler-utils
# Install air for hot-reloading
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Set PATH for binaries
ENV PATH=$PATH:/go/bin

CMD ["air"]
