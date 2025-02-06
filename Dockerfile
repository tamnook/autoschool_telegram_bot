# Используем официальный образ Go
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

# Собираем приложение
RUN go build -o autoschool_telegram_bot .

# Указываем команду для запуска приложения
CMD ["./autoschool_telegram_bot"]