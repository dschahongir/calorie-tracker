# ЭТАП 1: Сборка бинарника
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем приложение
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o calories-checker ./cmd/api/main.go

# ЭТАП 2: Минимальный рантайм
FROM alpine:latest

WORKDIR /app

# Копируем сертификаты для работы HTTPS (нужно для запросов к Gemini API)
RUN apk --no-cache add ca-certificates

# Копируем скомпилированный бинарник
COPY --from=builder /app/calories-checker .

EXPOSE 8080

CMD ["./calories-checker"]
