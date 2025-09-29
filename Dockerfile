FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git bash gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o password-manager ./main.go


FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/password-manager .

COPY .env .

EXPOSE 8000

CMD ["./password-manager"]
