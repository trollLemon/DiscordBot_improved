# syntax=docker/dockerfile:1

FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy environment file
COPY .env ./

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o bot main.go


FROM golang:1.23.1 AS tester
WORKDIR /app
COPY main.go ./
COPY Application ./Application
COPY Core ./Core
COPY util ./util
COPY go.mod go.sum ./

RUN go mod download

CMD ["go", "test", "-v", "./..."]

# Runtime Stage
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y \
    ffmpeg \
    pulseaudio \
    libasound2 \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY --from=builder /app/bot .
COPY --from=builder /app/.env .

EXPOSE 8080

ENTRYPOINT ["/bot"]