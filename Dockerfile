FROM golang:1.19 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# Running multiple stage is much less resource intense
FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./app"]
