# Используем официальный образ Golang в качестве базового
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем приложение
RUN go build -o grpc-server ./cmd/main.go

# Указываем порты, которые будут использоваться
EXPOSE 50051
EXPOSE 8080

# Копируем файл .env в контейнер
COPY .env.docker .env

# Устанавливаем переменные окружения из файла .env
ENV $(cat .env | xargs)

# Команда для запуска приложения
CMD ["./grpc-server"]