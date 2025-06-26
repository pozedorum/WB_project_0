FROM golang:1.24.4-alpine


WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git postgresql-client

# Копируем ВСЕ файлы проекта
COPY . .

RUN go mod download

EXPOSE 8081

CMD ["go", "run", "./cmd/server/main.go"]