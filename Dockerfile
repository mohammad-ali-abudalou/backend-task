# --- Build stage ---
FROM golang:1.25 AS build
WORKDIR /app

# Copy and download dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole project
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o server ./cmd


# --- Runtime stage ---
