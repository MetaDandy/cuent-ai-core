# ---------- build stage ----------
    FROM golang:1.24.1-alpine AS builder

    # Directorio de trabajo para compilación
    WORKDIR /src
    
    # Dependencias del sistema para compilación
    RUN apk add --no-cache git gcc musl-dev
    
    # Descarga de dependencias de Go
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copia código y compila para Linux amd64
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api
    
    # ---------- runtime stage ----------
    FROM alpine:3.20 AS runtime
    
    # Directorio de trabajo de la aplicación
    WORKDIR /app
    
    # Instala certificados y copia binario compilado
    RUN apk add --no-cache ca-certificates
    COPY --from=builder /src/app .
    
    # Permisos de ejecución
    RUN chmod +x ./app
    
    # Puerto que inyecta Render
    EXPOSE ${PORT}
    
    # Arranque de la aplicación escuchando en todas las interfaces
    CMD ["./app"]
    