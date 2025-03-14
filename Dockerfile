# syntax=docker/dockerfile:1

FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy environment file
COPY .env ./

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o bot main.go

# Runtime Stage
FROM archlinux:latest

RUN pacman -Syu --noconfirm && \
    pacman -S --noconfirm \
    ffmpeg \
    pulseaudio \
    alsa-lib \
    yt-dlp   \
    && pacman -Scc --noconfirm

WORKDIR /

COPY --from=builder /app/bot .
COPY --from=builder /app/.env .

EXPOSE 8080

ENTRYPOINT ["/bot"]
