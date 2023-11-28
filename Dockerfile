# Usa la imagen de Go como constructora para compilar la aplicación
FROM golang:1.18 AS builder

# Establece el directorio de trabajo en el contenedor
WORKDIR /app

# Copia los archivos de módulos Go
COPY go.mod ./
COPY go.sum ./

# Descarga las dependencias del módulo Go
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Construye la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Usa la imagen de Alpine para la imagen final, que es una imagen pequeña y segura
FROM alpine:latest

# Establece el directorio de trabajo en el contenedor
WORKDIR /

# Copia el binario compilado y el archivo .env desde el constructor
COPY --from=builder /app/main .
COPY .env .env

# Expone el puerto en el que se ejecutará la aplicación
EXPOSE 8080

# Ejecuta la aplicación
CMD ["./main"]