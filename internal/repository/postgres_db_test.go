package repository

import (
	"context"
	"log"
	"testing"

	"github.com/Egor-Tihonov/GRPC/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	Pool *pgxpool.Pool
)

type Service struct { // Service new
	rps Repository
}

func NewService(newRps Repository) *Service { // create
	return &Service{rps: newRps}
}

func TestCreate(t *testing.T) {
	testValidData := []model.Person{
		{
			Name:     "Ivan",
			Works:    true,
			Age:      19,
			Password: "ivan2001",
		},
		{
			Name:     "query2",
			Works:    true,
			Age:      19,
			Password: "1",
		},
	}
	testNoValidData := []model.Person{
		{
			Name:     "Ivan",
			Works:    false,
			Age:      18,
			Password: "3",
		},
		{
			Name:     "qwerty",
			Works:    true,
			Age:      -5,
			Password: "4",
		},
		{
			Name:     "qwerty1",
			Works:    false,
			Age:      250,
			Password: "250",
		},
	}
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, p := range testValidData {
		_, err := rps.rps.CreateUser(ctx, &p)
		require.NoError(t, err, "create error")
	}
	for _, p := range testNoValidData {
		_, err := rps.rps.CreateUser(ctx, &p)
		require.Error(t, err, "create error")
	}
}
func TestSelectAll(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := model.Person{
		ID:       "12",
		Name:     "Andrey",
		Works:    true,
		Age:      20,
		Password: "12",
	}

	users, err := rps.rps.GetAllUsers(ctx)
	require.NoError(t, err, "select all: problems with select all users")
	require.Equal(t, 2, len(users), "select all: the values are`t equals")

	_, err = Pool.Exec(ctx, "insert into persons(id,name,works,age,password) values($1,$2,$3,$4,$5)", &p.ID, &p.Name, &p.Works, &p.Age, &p.Password)
	require.NoError(t, err, "select all: insert error")
	users, err = rps.rps.GetAllUsers(ctx)
	if err != nil {
		defer log.Fatalf("error with select all: %v", err)
	}
	require.NotEqual(t, 5, len(users), "select all: the values are equals")
}

func TestSelectById(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	_, err := rps.rps.GetUserByID(ctx, "12")
	require.NoError(t, err, "select user by id: this id dont exist")
	_, err = rps.rps.GetUserByID(ctx, "20")
	require.Error(t, err, "select user by id: this id already exist")
	cancel()
}

func TestUpdate(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testValidData := []*model.Person{
		{
			Name:  "Masha",
			Works: true,
			Age:   19,
		},
		{
			Name:  "query21",
			Works: true,
			Age:   19,
		},
	}
	testNoValidData := []*model.Person{
		{
			Name:  "Egor",
			Works: false,
			Age:   120,
		},
		{
			Name:  "qwerty",
			Works: true,
			Age:   -5,
		},
		{
			Name:  "qwerty1",
			Works: false,
			Age:   250,
		},
	}
	for _, p := range testValidData {
		err := rps.rps.UpdateUser(ctx, "d57d1026-c79a-443d-9d81-714381a37a80", p)
		require.NoError(t, err, "update error")
	}
	for _, p := range testNoValidData {
		err := rps.rps.UpdateUser(ctx, "bb839db7-4be3-41a8-a53b-403ad26593ca", p)
		require.Error(t, err, "update error")
	}
	err := rps.rps.UpdateUser(ctx, "bb839db7-4be3-41a8-a53b-403ad26593ca", testValidData[0])
	require.Error(t, err, "update error")
}
func TestPRepository_UpdateAuth(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := rps.rps.UpdateAuth(ctx, "d57d1026-c79a-443d-9d81-714381a37a80",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTg0OTk0NzgsImp0aSI6IjNhZjYyMjY5LTAxZmYtNGM2YS04MmUwLTBhNjIwZTVlY2ZmZCIsInVzZXJuYW1lIjoiRWdvclRpaG9ub3YifQ.d4kAjYeGkObPF-kcm7TaFRducO7rsUjabu_8h-Sy8ZE")
	require.NoError(t, err, "thereis an error")
	err = rps.rps.UpdateAuth(ctx, "3",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTg0OTk0NzgsImp0aSI6IjNhZjYyMjY5LTAxZmYtNGM2YS04MmUwLTBhNjIwZTVlY2ZmZCIsInVzZXJuYW1lIjoiRWdvclRpaG9ub3YifQ.d4kAjYeGkObPF-kcm7TaFRducO7rsUjabu_8h-Sy8ZE")
	require.Error(t, err, "there isnt an error")
}
func TestPRepository_SelectByIdAuth(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := rps.rps.SelectByIDAuth(ctx, "d57d1026-c79a-443d-9d81-714381a37a80")
	require.NoError(t, err, "there is an error")
	_, err = rps.rps.SelectByIDAuth(ctx, "3")
	require.Error(t, err, "there isn`t an error")
}

func TestPRepository_Delete(t *testing.T) {
	rps := NewService(&PRepository{Pool: Pool})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := rps.rps.SelectByIDAuth(ctx, "d57d1026-c79a-443d-9d81-714381a37a80")
	require.NoError(t, err, "there is an error")
	_, err = rps.rps.SelectByIDAuth(ctx, "3")
	require.Error(t, err, "there isn`t an error")
}
