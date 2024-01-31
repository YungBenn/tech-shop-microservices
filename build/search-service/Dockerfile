FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o search-service cmd/search/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/search-service /app
COPY --from=builder /app/.env /app
CMD [ "./search-service" ]