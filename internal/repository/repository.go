// Package repository a
package repository

import (
	"github.com/Egor-Tihonov/GRPC/internal/model"

	"context"
)

// Repository transition to mongo or postgres db
type Repository interface {
	CreateUser(ctx context.Context, p *model.Person) (string, error)
	GetUserByID(ctx context.Context, idPerson string) (*model.Person, error)
	GetAllUsers(ctx context.Context) ([]*model.Person, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, per *model.Person) error
	SelectByIDAuth(ctx context.Context, id string) (model.Person, error)
	UpdateAuth(ctx context.Context, id string, refreshToken string) error
}
