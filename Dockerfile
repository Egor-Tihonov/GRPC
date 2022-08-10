FROM golang:latest

WORKDIR /GRPC

COPY . .

CMD ["./main"]