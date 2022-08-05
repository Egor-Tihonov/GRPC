package main

import (
	"awesomeProjectGRPC/internal/repository"
	"awesomeProjectGRPC/internal/server"
	pb "awesomeProjectGRPC/proto"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"net"
)

var poolP pgxpool.Pool

func main() {
	conn := DBConnection("postgres")
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
func DBConnection(dbname string) repository.Repository {
	switch dbname {
	case "postgres":
		poolP, err := pgxpool.Connect(context.Background() /*cfg.PostgresDbUrl */, "postgresql://postgres:123@localhost:5432/person")
		if err != nil {
			log.Errorf("bad connection with postgresql: %v", err)
			return nil
		}
		return &repository.PRepository{Pool: poolP}

	case "mongo":
		poolM, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
		if err != nil {
			log.Errorf("bad connection with mongoDb: %v", err)
			return nil
		}
		return &repository.MRepository{Pool: poolM}
	}
	return nil
}
