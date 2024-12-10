FROM golang:1.22.2-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o genres .

RUN chmod +x genres
