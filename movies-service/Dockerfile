FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY movies-service/go.mod movies-service/go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/movies-service github.com/jamescookdev/projeto-sipub-tech/movies-service/cmd


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/movies-service .

COPY data ./data


EXPOSE 50051

CMD ["./movies-service"]