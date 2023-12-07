# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init && swag fmt
RUN CGO_ENABLED=0 GOOS=linux go build -o moilerplate-app

# Stage 2: Create the final image
FROM alpine:latest

RUN mkdir logs/

COPY --from=builder /app/moilerplate-app /src/moilerplate-app

ENTRYPOINT ["/src/moilerplate-app", "server", "--mode=PRODUCTION"]