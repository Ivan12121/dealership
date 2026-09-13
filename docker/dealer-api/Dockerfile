FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dealer-api.bin .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/dealer-api.bin .

RUN chmod +x /root/dealer-api.bin

EXPOSE 8080
CMD ["./dealer-api.bin"]
