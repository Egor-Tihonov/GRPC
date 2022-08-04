package server

import (
	"awesomeProjectGRPC/internal/model"
	"awesomeProjectGRPC/internal/repository"
	pb "awesomeProjectGRPC/proto"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Server struct {
	pb.UnimplementedCRUDServer
	rps repository.Repository
}

var (
	accessTokenWorkTime  = time.Now().Add(time.Minute * 5).Unix()
	refreshTokenWorkTime = time.Now().Add(time.Hour * 3).Unix()
)
var JwtKey = []byte("super-key")

func NewServer(pool repository.Repository) *Server {
	return &Server{rps: pool}
}
func (s *Server) Authentication(ctx context.Context, request *pb.AuthenticationRequest) (*pb.Response, error) {
	authUser, err := s.rps.GetUserByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	incoming := []byte(request.Password)
	existing := []byte(authUser.Password)
	err = bcrypt.CompareHashAndPassword(existing, incoming) // check passwords
	if err != nil {
		return nil, err
	}
	authUser.Password = request.Password

	return nil, nil
}
func (s *Server) Registration(ctx context.Context, request *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	hashPassword, err := hashingPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("server: error while hashing password, %e", err)
	}
	request.Password = hashPassword
	p := model.Person{
		Name:     request.Name,
		Works:    request.Works,
		Age:      request.Age,
		Password: request.Password,
	}
	newID, err := s.rps.CreateUser(ctx, p)
	if err != nil {
		return nil, err
	}
	return &pb.RegistrationResponse{Id: newID}, nil
}

func (s *Server) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
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
func (s *Server) RefreshToken(ctx context.Context, refreshTokenString string) (accessTokenStr, refreshTokenStr string, err error) { // refresh our tokens
	refreshToken, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	}) // parse it into string format
	if err != nil {
		log.Errorf("service: can't parse refresh token - %e", err)
		return "", "", err
	}
	if !refreshToken.Valid {
		return "", "", fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		return "", "", fmt.Errorf("service: error while parsing claims, ID couldnt be empty")
	}
	person, err := s.rps.SelectByIDAuth(ctx, userUUID.(string))
	if err != nil {
		return "", "", fmt.Errorf("service: token refresh failed - %e", err)
	}
	if refreshTokenString != person.RefreshToken {
		return "", "", fmt.Errorf("service: invalid refresh token")
	}
	return s.CreateJWT(ctx, s.rps, &person)
}

// CreateJWT create jwt tokens
func (s *Server) CreateJWT(ctx context.Context, rps repository.Repository, person *model.Person) (accessTokenStr, refreshTokenStr string, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)         // encrypt access token by SigningMethodHS256 method
	claimsA := accessToken.Claims.(jwt.MapClaims)          // fill access-token`s claims
	claimsA["exp"] = accessTokenWorkTime                   // work time
	claimsA["username"] = person.Name                      // payload
	accessTokenStr, err = accessToken.SignedString(JwtKey) // convert token to string format
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claimsR := refreshToken.Claims.(jwt.MapClaims)
	claimsR["username"] = person.Name
	claimsR["exp"] = refreshTokenWorkTime
	claimsR["jti"] = person.ID
	refreshTokenStr, err = refreshToken.SignedString(JwtKey)
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	err = rps.UpdateAuth(ctx, person.ID, refreshTokenStr) // add into user refresh token
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	return
}

// UpdateUserAuth update auth user, add token
func (s *Server) UpdateUserAuth(ctx context.Context, id, refreshToken string) error {
	return s.rps.UpdateAuth(ctx, id, refreshToken)
}

func hashingPassword(password string) (string, error) {
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword)
	return hashPassword, nil
}
