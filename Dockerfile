# syntax=docker/dockerfile:1

FROM golang:1.23.1 AS builder

WORKDIR /app

COPY cmd ./cmd 
COPY internal ./internal/  
COPY go.mod go.sum ./

RUN CGO_ENABLED=1 GOOS=linux go build -o bot ./cmd/*


FROM golang:1.23.1 AS tester
WORKDIR /app
COPY cmd ./
COPY internal ./ 
COPY go.mod go.sum ./

RUN go mod download

CMD ["go", "test", "-v", "./..."]

# Runtime Stage
FROM debian:bookworm-slim

WORKDIR /
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=builder /app/bot .

EXPOSE 8080

CMD ["/bot", "--pretty-print"]
