FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/

RUN go build -o app cmd/app/main.go

FROM alpine:latest

COPY --from=builder /app/app /app/app

CMD ["/app/app"]
