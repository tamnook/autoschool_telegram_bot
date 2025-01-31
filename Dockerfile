# Используем официальный образ Go
FROM golang:1.21-alpine

RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

RUN git pull .

RUN go mod tidy

# Собираем приложение
RUN go build -o autoschool_telegram_bot .

# Указываем команду для запуска приложения
CMD ["./autoschool_telegram_bot"]