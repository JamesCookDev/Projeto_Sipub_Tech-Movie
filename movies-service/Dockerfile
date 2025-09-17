# Estágio 1: 'builder' - Apenas compila a aplicação
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/movies-service ./cmd/main.go


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/movies-service .

EXPOSE 50051

CMD ["./movies-service"]