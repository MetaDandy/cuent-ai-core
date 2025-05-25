# ---------- build stage ----------
    FROM golang:1.24.1-alpine AS builder
    WORKDIR /src
    RUN apk add --no-cache git gcc musl-dev
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api
    
    # ---------- runtime stage ----------
    FROM alpine:3.20
    WORKDIR /app
    RUN apk add --no-cache ca-certificates
    COPY --from=builder /src/app .
    RUN chmod +x ./app
    ENV PORT 8000
    EXPOSE ${PORT}
    CMD ["./app"]
    