// Package repository : file contains operations with MongoDB
package repository

import (
	"context"
	"fmt"

	"github.com/Egor-Tihonov/GRPC/internal/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MRepository create connection with MongoD
type MRepository struct {
	Pool *mongo.Client
}

// CreateUser add new user to db
func (m *MRepository) CreateUser(ctx context.Context, person *model.Person) (string, error) {
	if person.Age < 0 || person.Age > 180 {
		return "", fmt.Errorf("mongo repository: error with create, age must be more then 0 and less then 180")
	}
	newID := uuid.New().String()
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: newID},
		{Key: "name", Value: person.Name},
		{Key: "works", Value: person.Works},
		{Key: "age", Value: person.Age},
		{Key: "password", Value: person.Password},
		{Key: "refreshtoken", Value: person.RefreshToken},
	})
	if err != nil {
		return "", fmt.Errorf("mongo: unable to create new user: %v", err)
	}
	return newID, nil
}

// UpdateUser update exist user
func (m *MRepository) UpdateUser(ctx context.Context, id string, person *model.Person) error {
	if person.Age < 0 || person.Age > 180 {
		return fmt.Errorf("mongo repository: error with create, person`s age must be more then 0 and less then 180")
	}
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: person.Name},
		{Key: "works", Value: person.Works},
		{Key: "age", Value: person.Age},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update user %v", err)
	}
	return nil
}

// UpdateAuth add user refresh token
func (m *MRepository) UpdateAuth(ctx context.Context, id, refreshToken string) error {
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.UpdateOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "refreshtoken", Value: refreshToken},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update user %v", err)
	}
	return nil
}

// GetAllUsers take all users from db
func (m *MRepository) GetAllUsers(ctx context.Context) ([]*model.Person, error) {
	var users []*model.Person
	collection := m.Pool.Database("person").Collection("person")
	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongo: unable to select all users %v", err)
	}
	for c.Next(ctx) {
		user := model.Person{}
		err := c.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// DeleteUser user from db
func (m *MRepository) DeleteUser(ctx context.Context, id string) error {
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.DeleteOne(ctx, bson.D{primitive.E{Key: "id", Value: id}})
	if err != nil {
		return fmt.Errorf("mongo: unable to delete user, %v", err)
	}
	return nil
}

// GetUserByID select exist user from db by his id
func (m *MRepository) GetUserByID(ctx context.Context, id string) (*model.Person, error) {
	user := model.Person{}
	collection := m.Pool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}).Decode(&user)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

// SelectByIDAuth take from user his refresh token
func (m *MRepository) SelectByIDAuth(ctx context.Context, id string) (model.Person, error) {
	user := model.Person{}
	collection := m.Pool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: id},
		{Key: "name", Value: 0},
		{Key: "works", Value: 0},
		{Key: "age", Value: 0},
		{Key: "password", Value: 0},
	}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
