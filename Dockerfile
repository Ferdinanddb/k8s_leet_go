FROM golang:1.21.6-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/backend

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/backend .

ENV DB_HOST="127.0.0.1"
ENV DB_USER="postgres"
ARG DB_PASSWORD
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_NAME="test"
ENV DB_PORT="5432"
ENV TOKEN_TTL="2000"
ARG JWT_PRIVATE_KEY
ENV JWT_PRIVATE_KEY=$JWT_PRIVATE_KEY

CMD ["./backend"]