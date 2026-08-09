FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /barberflow \
    ./cmd/barberflow

FROM alpine:3.22

COPY --from=builder /barberflow /app/barberflow

CMD ["/app/barberflow"]