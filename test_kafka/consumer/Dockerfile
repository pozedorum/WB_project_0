FROM golang:1.21-alpine

WORKDIR /app

# Копируем только файлы модуля сначала для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY *.go ./

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer .

CMD ["./consumer"]