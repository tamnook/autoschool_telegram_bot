# Используем официальный образ Go
FROM golang:latest-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

# Собираем приложение
RUN go build -o autoschool_telegram_bot .

# Указываем команду для запуска приложения
CMD ["./autoschool_telegram_bot"]