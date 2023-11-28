# Stage 1: Build the application
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o moilerplate-app

# Stage 2: Create the final image
FROM alpine:latest

RUN mkdir logs/

COPY --from=builder /app/moilerplate-app /src/moilerplate-app

EXPOSE 80

ENTRYPOINT ["/src/moilerplate-app server --mode=PRODUCTION"]
