package main

import (
	"awesomeProjectGRPC/internal/repository"
	"awesomeProjectGRPC/internal/server"
	pb "awesomeProjectGRPC/proto"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
)

var poolP pgxpool.Pool

func main() {
	conn := DBConnection()
	defer poolP.Close()
	ns := grpc.NewServer()
	srv := server.NewServer(conn)
	pb.RegisterCRUDServer(ns, srv)
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error while listening port: %e", err)
	}
	log.Print("success listen server, port 50051")
	if err = ns.Serve(listen); err != nil {
		log.Fatalf("error while listening server: %e", err)
	}

}
func DBConnection() repository.Repository {
	poolP, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
	if err != nil {
		log.Fatalf("bad connection with postgresql: %v", err)
		return nil
	}
	log.Print("success connect to postgres db")
	return &repository.PRepository{Pool: poolP}
}
