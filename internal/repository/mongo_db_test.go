package repository

import (
	"log"
	"os"
	"testing"

	"github.com/Egor-Tihonov/GRPC/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var (
	PoolM *mongo.Client
)

type ServiceM struct { //service new
	rps Repository
}

func NewServiceM(NewRps Repository) *ServiceM { //create
	return &ServiceM{rps: NewRps}
}
func TestMain(m *testing.M) {
	pool, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatalf("Bad connection: %v", err)
	}
	PoolM = pool
	run := m.Run()
	os.Exit(run)
}
func TestCreateMongo(t *testing.T) {
	testValidData := []model.Person{
		{
			Name:     "Ivan",
			Works:    true,
			Age:      19,
			Password: "0",
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
	}
	rps := NewServiceM(&MRepository{Pool: PoolM})
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
func TestSelectAllMongo(t *testing.T) {
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	users, err := rps.rps.GetAllUsers(ctx)
	require.NoError(t, err, "select all: problems with select all users")
	require.Equal(t, 2, len(users), "select all: the values are`t equals")
	collection := PoolM.Database("person").Collection("person")
	_, err = collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: "1223"},
		{Key: "name", Value: "Sasha"},
		{Key: "works", Value: true},
		{Key: "age", Value: 18},
		{Key: "password", Value: "sheisverybeatiful"},
		{Key: "refreshtoken", Value: ""},
	})
	require.NoError(t, err, "select all: insert error")
	users, _ = rps.rps.GetAllUsers(ctx)
	require.Equal(t, 3, len(users), "select all: the values are`t equals")
	require.NotEqual(t, 4, len(users), "select all: the values are equals")
}
func TestSelectByIdMongo(t *testing.T) {
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	_, err := rps.rps.GetUserByID(ctx, "1223")
	require.NoError(t, err, "select user by id: this id dont exist")
	_, err = rps.rps.GetUserByID(ctx, "20")
	require.Error(t, err, "select user by id: this id already exist")
	cancel()
}
func TestUpdateMongo(t *testing.T) {
	rps := NewServiceM(&MRepository{Pool: PoolM})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testValidData := []model.Person{
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
	testNoValidData := []model.Person{
		{
			Name:  "Ivan",
			Works: false,
			Age:   12,
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
		err := rps.rps.UpdateUser(ctx, "1223", &p)
		require.NoError(t, err, "update error")
	}
	for _, p := range testNoValidData {
		err := rps.rps.UpdateUser(ctx, "1223", &p)
		require.Error(t, err, "update error")
	}
}
