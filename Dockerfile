# =========================
# Stage 1: Build
# =========================

FROM golang:1.25 AS builder

WORKDIR /app

# Copy Go Modules Manifests :
COPY go.mod go.sum ./

# Download Dependencies :
RUN go mod download

# Copy The Source Code :
COPY . .

# Build The Go Binary  :
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o backend-task ./cmd/app


# =========================
# Stage 2: Runtime
# =========================

FROM alpine:3.19

WORKDIR /root/

# Copy The Binary From Builder Stage :
COPY --from=builder /app/backend-task .

# Expose The API Port :
EXPOSE 8080

# Run The Application :
CMD ["./backend-task"]
