FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o order-service cmd/order/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/order-service /app
COPY --from=builder /app/.env /app
CMD [ "./order-service" ]