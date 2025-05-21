# Dockerfile
# Gunakan image base Go yang lebih kecil
FROM golang:1.22-alpine AS builder

# Atur working directory
WORKDIR /app

# Salin file go.mod dan go.sum untuk dependensi caching
COPY go.mod ./
COPY go.sum ./

# Unduh dependensi
RUN go mod download

# Salin source code aplikasi
COPY . .

# Bangun aplikasi
# CGO_ENABLED=0 diperlukan untuk static linking saat menggunakan Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Gunakan image base yang sangat kecil untuk binary akhir
FROM alpine:latest

WORKDIR /root/

# Salin binary yang sudah dibangun dari builder stage
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi
EXPOSE 9688

# Jalankan aplikasi
CMD ["./main"]