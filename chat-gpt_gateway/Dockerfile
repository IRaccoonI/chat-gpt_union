# Определение базового образа
FROM golang:latest

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копирование файлов с кодом в контейнер
COPY . .

# Сборка приложения
RUN go build -o app

EXPOSE 3000

# Определение команды запуска приложения при запуске контейнера
CMD ["./app"]
