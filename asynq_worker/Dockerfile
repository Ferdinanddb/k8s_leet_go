FROM golang:1.21.6-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/asynq_worker

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/asynq_worker .

ENV REDIS_HOST=$REDIS_HOST
ENV REDIS_PORT=$REDIS_PORT
ENV REDIS_PASSWORD=$REDIS_PASSWORD

CMD ["./asynq_worker"]