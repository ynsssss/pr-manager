FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr-manager ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pr-manager .

EXPOSE 8080

CMD ["./pr-manager"]

