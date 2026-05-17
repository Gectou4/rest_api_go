FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates mariadb-client

WORKDIR /app

COPY --from=builder /app/server .
COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:3306", "--", "./server"]
