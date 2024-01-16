FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o tech-shop-gateway ./cmd/tech-shop-gateway/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/tech-shop-gateway /app
COPY --from=builder /app/.env /app
CMD [ "./tech-shop-gateway" ]