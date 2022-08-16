package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Egor-Tihonov/GRPC/internal/model"
	"github.com/Egor-Tihonov/GRPC/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// NewServer create new server connection
var (
	Pool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://postgres:123@localhost:5432/person")
	if err != nil {
		log.Fatalf("Bad connection: %v", err)
	}
	Pool = pool
	run := m.Run()
	os.Exit(run)
}

func TestAuthentication(t *testing.T) {
	se := NewService(&repository.PRepository{Pool: Pool}, []byte("super-key"))
	_, _, err := se.Authentication(context.Background(), "c651c18a-3ff8-436f-a543-74df9f5f5706", "tujh2004")
	require.NoError(t, err, "")
	_, _, err = se.Authentication(context.Background(), "c651c18a-3ff8-436f-a543-74df9f5f5706", "2")
	require.Error(t, err, "Password is available for this person")
	_, _, err = se.Authentication(context.Background(), "c651c18a-3ff8-436f-a543-74df9f5f5701", "tujh2004")
	require.Error(t, err, "Password is available for this person")
}

func TestCreateUser(t *testing.T) {
	se := NewService(&repository.PRepository{Pool: Pool}, []byte("super-key"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testValidData := model.Person{Name: "Vanua", Works: true, Age: 19, Password: "worker"}
	testNoValidData := model.Person{Name: "Mani", Works: true, Age: 19, Password: "worker"}
	_, err := se.CreateUser(ctx, &testValidData)
	require.NoError(t, err, "cannot register this user")
	_, err = se.CreateUser(ctx, &testNoValidData)
	require.Error(t, err, "this user can be register")

}

func TestHashPassword(t *testing.T) {
	pass, err := hashingPassword("1234567890")
	require.NoError(t, err, "")
	_, err = hashingPassword("12")
	require.Error(t, err, "")
	incoming := []byte("1234567890")
	existing := []byte(pass)
	err = bcrypt.CompareHashAndPassword(existing, incoming)
	require.NoError(t, err, "")
	incoming = []byte("1234567891")
	existing = []byte(pass)
	err = bcrypt.CompareHashAndPassword(existing, incoming)
	require.Error(t, err, "")
}

func TestVerifyAccessToken(t *testing.T) {
	se := NewService(&repository.PRepository{Pool: Pool}, []byte("super-key"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	access, _, _ := se.Authentication(ctx, "c651c18a-3ff8-436f-a543-74df9f5f5706", "tujh2004")

	err := se.Verify(access)
	require.NoError(t, err)
	err = se.Verify("123123123123123123123123132312313231231323132312313231231")
	require.Error(t, err)
}

func TestCreateJWT(t *testing.T) {
	testUser := model.Person{
		ID:       "c651c18a-3ff8-436f-a543-74df9f5f5706",
		Name:     "Egor Tihonov",
		Works:    true,
		Age:      18,
		Password: "tujh2004",
	}
	testUserNoValidate := model.Person{
		ID:       "0",
		Name:     "Egor Tihonov",
		Works:    true,
		Age:      120,
		Password: "tujh2004",
	}
	s := NewService(&repository.PRepository{Pool: Pool}, []byte("super-key"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, _, err := s.CreateJWT(ctx, s.rps, &testUser)
	require.NoError(t, err, "cannot create tokens")
	_, _, err = s.CreateJWT(ctx, s.rps, &testUserNoValidate)
	require.Error(t, err, "tokens create")
}

func TestRefreshToken(t *testing.T) {
	s := NewService(&repository.PRepository{Pool: Pool}, []byte("super-key"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, _, err := s.RefreshTokens(ctx, "token")
	require.NoError(t, err, "cannot refresh your tokens")
	_, _, err = s.RefreshTokens(ctx, "<false token>")
	require.Error(t, err, "can refresh your tokens")
	_, _, err = s.RefreshTokens(ctx, "old token")
	require.Error(t, err, "token already valid")
}
