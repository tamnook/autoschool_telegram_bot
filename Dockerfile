# Используем официальный образ Go
FROM golang:latest-alpine

RUN apk add --no-cache git

COPY . /app

WORKDIR /app

RUN git pull .

RUN go mod tidy

# Собираем приложение
RUN go build -o autoschool_telegram_bot .

# Указываем команду для запуска приложения
CMD ["./autoschool_telegram_bot"]