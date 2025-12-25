FROM golang:1.25.0-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o build ./main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/build .

CMD ["./build"]

