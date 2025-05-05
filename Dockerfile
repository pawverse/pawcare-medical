FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:latest

ARG HTTP_PORT=80

ENV HTTP_PORT=${HTTP_PORT}

WORKDIR /app
COPY --from=builder /app/server /app

EXPOSE ${HTTP_PORT}

USER nobody:nobody

ENTRYPOINT ["/app/server"]
