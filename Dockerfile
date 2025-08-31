# --- Build stage ---
FROM golang:1.25 AS build
WORKDIR /app

# Copy and download dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole project
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false-o backend-task server ./cmd/app

# Expose API port
EXPOSE 8080

# --- Runtime stage ---

# ============================
# Stage 2: Run
# ============================
FROM alpine:3.19

WORKDIR /root/

# Expose API port
EXPOSE 8080

# Run the app
CMD ["./backend-task"]