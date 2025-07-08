FROM golang:1.23.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o qr-quest ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/qr-quest .

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/internal/handlers/ComicSansMS.ttf /app/fonts/ComicSansMS.ttf

EXPOSE 8080
CMD ["./qr-quest"]
