# ---------- build stage ----------
    FROM golang:1.24.1-alpine AS builder
    WORKDIR /src
    
    # dependencias de compilación
    RUN apk add --no-cache git gcc musl-dev
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    # compilamos explícitamente para linux/amd64
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o app ./cmd/api
    
    # ---------- runtime stage ----------
    FROM alpine:3.20
    WORKDIR /app
    
    # copiamos el binario ya compilado
    COPY --from=builder /src/app .
    
    # permisos de ejecución
    RUN chmod +x app
    
    # Render inyecta PORT automáticamente
    EXPOSE ${PORT}
    
    # arrancamos escuchando en 0.0.0.0:${PORT}
    CMD ["./app"]
    