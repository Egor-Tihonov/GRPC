package server

import (
	"awesomeProjectGRPC/internal/model"
	"awesomeProjectGRPC/internal/repository"
	pb "awesomeProjectGRPC/proto"
	"context"
	"time"
)

type Server struct {
	pb.UnimplementedCRUDServer
	rps repository.Repository
}

var JwtKey = []byte("super-key")

var (
	AccessTokenWorkTime  = time.Now().Add(time.Minute * 5).Unix()
	RefreshTokenWorkTime = time.Now().Add(time.Hour * 3).Unix()
)

func NewServer(pool repository.Repository) *Server {
	return &Server{rps: pool}
}

func (s *Server) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := Verify(request.AccessToken); err != nil {
		return nil, err
	}
	idPerson := request.GetId()
	personDB, err := s.rps.GetUserByID(ctx, idPerson)
	if err != nil {
		return nil, err
	}
	personProto := &pb.GetUserResponse{
		Person: &pb.Person{
			Id:       personDB.ID,
			Name:     personDB.Name,
			Age:      personDB.Age,
			Works:    personDB.Works,
			Password: personDB.Password,
		},
	}
	return personProto, nil
}

func (s *Server) GetAllUsers(_ *pb.GetAllUsersRequest, stream pb.CRUD_GetAllUsersServer) error {
	persons, err := s.rps.GetAllUsers(context.Background())
	var personProto pb.Person
	if err != nil {
		return err
	}
	for _, person := range persons {
		personProto.Id = person.ID
		personProto.Name = person.Name
		personProto.Age = person.Age
		personProto.Works = person.Works
		err = stream.Send(&personProto)
		if err != nil {
			return err
		}

	}
	return err
}
func (s *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.Response, error) {
	err := s.rps.DeleteUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}

func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.Response, error) {
	if err := Verify(request.AccessToken); err != nil {
		return nil, err
	}
	if err := Verify(request.AccessToken); err != nil {
		return nil, err
	}
	user := &model.Person{
		Name:  request.Person.Name,
		Works: request.Person.Works,
		Age:   request.Person.Age,
	}
	err := s.rps.UpdateUser(ctx, request.Id, user)
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}
