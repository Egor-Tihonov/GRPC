// Package server p
package server

import (
	"awesomeProjectGRPC/internal/model"
	"awesomeProjectGRPC/internal/service"
	pb "awesomeProjectGRPC/proto"
	"context"
)

// Server struct
type Server struct {
	pb.UnimplementedCRUDServer
	se *service.Service
}

// NewServer create new server connection
func NewServer(serv *service.Service) *Server {
	return &Server{se: serv}
}

// GetUser get user by id from db
func (s *Server) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	accessToken := request.GetAccessToken()
	if err := s.se.Verify(accessToken); err != nil {
		return nil, err
	}
	idPerson := request.GetId()
	personDB, err := s.se.GetUser(ctx, idPerson)
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

// GetAllUsers get all users from db
func (s *Server) GetAllUsers(_ *pb.GetAllUsersRequest, stream pb.CRUD_GetAllUsersServer) error {
	persons, err := s.se.GetAllUsers()
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

// DeleteUser delete user by id
func (s *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.Response, error) {
	idUser := request.GetId()
	err := s.se.DeleteUser(ctx, idUser)
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}

// UpdateUser update user with new parameters
func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.Response, error) {
	accessToken := request.GetAccessToken()
	if err := s.se.Verify(accessToken); err != nil {
		return nil, err
	}
	user := &model.Person{
		Name:  request.Person.Name,
		Works: request.Person.Works,
		Age:   request.Person.Age,
	}
	idUser := request.GetId()
	err := s.se.UpdateUser(ctx, idUser, user)
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil
}
