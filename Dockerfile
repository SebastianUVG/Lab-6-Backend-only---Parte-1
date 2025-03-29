FROM golang:1.24.1-alpine AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY LaLigaTracker.html .

ENV DB_HOST=db
ENV DB_USER=POSTGRES
ENV DB_PASSWORD=Admin123
ENV DB_NAME=laligadb

EXPOSE 8080

CMD ["./main"]
