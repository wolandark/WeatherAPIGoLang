# FROM golang:1.21-alpine AS builder
# FROM golang:1.24.4-alpine AS builder

FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# May god strike alpine and its non-gnu bs down
# RUN go build -o weather-api ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/weather-api .
COPY .env .

EXPOSE 8080

CMD ["./weather-api"]

