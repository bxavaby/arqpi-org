FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/api .
COPY data ./data

RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./api"]
