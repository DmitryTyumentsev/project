# Используем официальный образ Go
FROM golang:1.21-alpine

# Устанавливаем зависимости (если нужны)
RUN apk add --no-cache git

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бота
RUN go build -o bot

# Запускаем
CMD ["./bot"]