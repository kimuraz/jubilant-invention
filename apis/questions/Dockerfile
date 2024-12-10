FROM golang:1.22.2-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o questions .

RUN chmod +x questions