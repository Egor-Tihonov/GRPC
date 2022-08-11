FROM golang:latest

WORKDIR /

COPY . .

CMD ["./main"]