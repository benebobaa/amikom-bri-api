# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# Set execute permissions for scripts
RUN chmod +x wait-for.sh
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY  app.env .
COPY  start.sh .
COPY  wait-for.sh .
COPY  db/migrations ./migrations
COPY public/expenses_plans ./public/expenses_plans
COPY public/transaction_pdf ./public/transaction_pdf


EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["sh", "/app/start.sh"]
