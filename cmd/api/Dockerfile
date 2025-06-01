FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/api /app/config_container.yaml .
COPY --from=builder /app/static/ ./static/

ENV CONFIG_PATH=/root/config_container.yaml

# ENTRYPOINT ["./api"]

CMD ["./api"]