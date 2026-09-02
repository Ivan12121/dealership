FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY dealer-api/ .

RUN go mod init dealer-api || true
RUN go mod tidy
RUN go build -o dealer-api-bin ./main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/dealer-api-bin .

EXPOSE 8080

CMD ["./dealer-api-bin"]