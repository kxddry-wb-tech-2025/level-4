FROM golang:1.25.0 AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/

RUN go build -o calendar ./cmd/calendar/main.go

FROM golang:1.25.0
WORKDIR /app

COPY --from=builder /app/calendar .

CMD ["./calendar"]