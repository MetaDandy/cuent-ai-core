# ---------- build stage ----------
    FROM golang:1.24.1-alpine AS builder

    WORKDIR /src
    
    # dependencias del sistema
    RUN apk add --no-cache git gcc musl-dev
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    RUN CGO_ENABLED=0 go build -o app ./cmd/api
    
    # ---------- runtime stage ----------
    FROM alpine:3.20
    
    WORKDIR /app
    COPY --from=builder /src/app .
    
    CMD ["./app"]
    