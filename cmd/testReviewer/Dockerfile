FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/testReviewer
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/testReviewer

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/testReviewer .

ENTRYPOINT ["./testReviewer"]