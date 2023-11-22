# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN mkdir logs

RUN go install github.com/cosmtrek/air@latest
RUN go mod tidy

CMD ["air", "-c", ".air.toml"]
