# Используем базовый образ Go. Обновляем версию до 1.24, чтобы соответствовать требованиям air.
# Этот образ содержит все необходимые инструменты для компиляции.
FROM golang:1.24-alpine

# Устанавливаем рабочую директорию внутри контейнера.
WORKDIR /app

# Устанавливаем 'air' для live-reloading.
# Эта утилита будет следить за изменениями в файлах и перезапускать сервер.
# Мы устанавливаем ее глобально в контейнере.
RUN go install github.com/air-verse/air@latest

# Копируем файлы go.mod и go.sum, чтобы Docker мог кешировать зависимости.
COPY go.mod go.sum ./

# Загружаем зависимости проекта.
RUN go mod download

# Копируем остальные файлы проекта в рабочую директорию контейнера.
# При использовании docker-compose с volumes этот шаг важен для первого билда,
# а затем локальная папка будет "зеркалироваться" в контейнер.
COPY . .

# Открываем порт 8080, чтобы можно было получить доступ к приложению.
EXPOSE 8080

# Команда по умолчанию для запуска контейнера.
# Запускаем 'air', который, в свою очередь, скомпилирует и запустит ваше приложение
# согласно настройкам из файла .air.toml.
CMD ["air"] 