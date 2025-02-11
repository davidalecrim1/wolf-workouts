FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM builder AS dev
WORKDIR /app
COPY .air.toml ./

VOLUME ["/app"]

CMD ["air", "-c", ".air.toml"]
