FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:latest as runner

ENV GO_ENV=production

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml .

EXPOSE 8080

CMD ["./main"]
