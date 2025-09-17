FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./

COPY movies-service/go.mod movies-service/go.sum ./movies-service/
COPY api/go.mod api/go.sum ./api/

RUN go work sync
RUN go mod download

COPY ./movies-service ./movies-service/
COPY ./api ./api/

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./api/


FROM alpine:latest
COPY --from=builder api-gateway /api-gateway
EXPOSE 8080
CMD ["/api-gateway"]