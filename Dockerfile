# --- Build Stage ---
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o smart-travel-advisor ./cmd/server/main.go

# --- Runtime Stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/smart-travel-advisor .

EXPOSE 8080

CMD ["./smart-travel-advisor"]
