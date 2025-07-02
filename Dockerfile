FROM golang:1.24.2 AS builder

WORKDIR /cmd

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

RUN ls -la /cmd

FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /cmd/server .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml
USER appuser

EXPOSE 8080

CMD ["./server"]