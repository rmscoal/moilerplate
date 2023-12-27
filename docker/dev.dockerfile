# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN mkdir logs

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/cosmtrek/air@latest
RUN go mod tidy

CMD ["air", "-c", ".air.toml"]
