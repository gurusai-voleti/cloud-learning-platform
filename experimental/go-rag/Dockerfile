FROM golang:1.23.2-alpine

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache \
    git \
    build-base \
    poppler-utils \  # Required for PDF processing
    libreoffice     # Required for Office document processing

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o document-processor ./main.go

ENTRYPOINT ["/app/document-processor"]
