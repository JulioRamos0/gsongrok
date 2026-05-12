# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gsongrok .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/gsongrok .
COPY --from=builder /app/public ./public

# Default environment variables
ENV PORT=8080
ENV PATHS_JSON_PATH=/data/paths.json

EXPOSE 8080

CMD ["./gsongrok"]
