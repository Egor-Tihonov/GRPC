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
)

//Authentication login
func (s *Server) Authentication(ctx context.Context, request *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
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
	accessToken, refreshToken, err := s.CreateJWT(ctx, s.rps, authUser)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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

func (s *Server) RefreshMyTokens(ctx context.Context, refreshTokenString *pb.RefreshTokensRequest) (*pb.RefreshTokensResponse, error) { // refresh our tokens
	refreshToken, err := jwt.Parse(refreshTokenString.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	}) // parse it into string format
	if err != nil {
		log.Errorf("service: can't parse refresh token - %e", err)
		return nil, err
	}
	if !refreshToken.Valid {
		return nil, fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		return nil, fmt.Errorf("service: error while parsing claims, ID couldnt be empty")
	}
	person, err := s.rps.SelectByIDAuth(ctx, userUUID.(string))
	if err != nil {
		return nil, fmt.Errorf("service: token refresh failed - %e", err)
	}
	if refreshTokenString.RefreshToken != person.RefreshToken {
		return nil, fmt.Errorf("service: invalid refresh token")
	}
	newAccessToken, newRefreshToken, err := s.CreateJWT(ctx, s.rps, &person)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokensResponse{
		RefreshToken: newRefreshToken,
		AccessToken:  newAccessToken,
	}, nil
}

// CreateJWT create jwt tokens
func (s *Server) CreateJWT(ctx context.Context, rps repository.Repository, person *model.Person) (accessTokenStr, refreshTokenStr string, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)         // encrypt access token by SigningMethodHS256 method
	claimsA := accessToken.Claims.(jwt.MapClaims)          // fill access-token`s claims
	claimsA["exp"] = AccessTokenWorkTime                   // work time
	claimsA["username"] = person.Name                      // payload
	accessTokenStr, err = accessToken.SignedString(JwtKey) // convert token to string format
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claimsR := refreshToken.Claims.(jwt.MapClaims)
	claimsR["username"] = person.Name
	claimsR["exp"] = RefreshTokenWorkTime
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
func (s *Server) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.Response, error) {
	err := Verify(request.AccessToken)
	if err != nil {
		return nil, err
	}
	err = s.UpdateUserAuth(ctx, request.Id, "")
	if err != nil {
		return nil, err
	}
	return new(pb.Response), nil

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

func Verify(accessTokenString string) error {
	accessToken, err := jwt.Parse(accessTokenString, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		log.Errorf("service: can't parse refresh token - ", err)
		return err
	}
	if !accessToken.Valid {
		return fmt.Errorf("service: expired refresh token")
	}
	return nil
}
