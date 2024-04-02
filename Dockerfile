# Builder
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY cmd ./cmd

COPY internal ./internal

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

RUN chmod +x ./main

# Production 
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]