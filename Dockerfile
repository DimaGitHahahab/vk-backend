FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o vk-api ./cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /app/.env .

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/vk-api .

ENTRYPOINT ["./vk-api"]
