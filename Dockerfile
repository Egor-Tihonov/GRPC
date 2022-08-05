FROM golang:latest

WORKDIR /awesomeProjectGRPC

COPY ./ /awesomeProjectGRPC

RUN go mod download

ENTRYPOINT go run main.go