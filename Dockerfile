# Используем официальный образ Go
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

# Собираем приложение
RUN go build -o -v autoschool_telegram_bot .

RUN ls -la /app

RUN chmod +x autoschool_telegram_bot .

# Указываем команду для запуска приложения
CMD ["./autoschool_telegram_bot"]