package repository

import (
	"github.com/Egor-Tihonov/GRPC/internal/model"

	"context"

	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

// PRepository p
type PRepository struct {
	Pool *pgxpool.Pool
}

// CreateUser add user to db
func (p *PRepository) CreateUser(ctx context.Context, person model.Person) (string, error) {
	newID := uuid.New().String()
	_, err := p.Pool.Exec(ctx, "insert into persons(id,name,works,age,password) values($1,$2,$3,$4,$5)",
		newID, &person.Name, &person.Works, &person.Age, &person.Password)
	if err != nil {
		log.Errorf("database error with create user: %v", err)
		return "", err
	}
	return newID, nil
}

// GetUserByID select user by id
func (p *PRepository) GetUserByID(ctx context.Context, idPerson string) (*model.Person, error) {
	u := model.Person{}
	err := p.Pool.QueryRow(ctx, "select id,name,works,age,password from persons where id=$1", idPerson).Scan(
		&u.ID, &u.Name, &u.Works, &u.Age, &u.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &model.Person{}, fmt.Errorf("user with this id doesnt exist: %v", err)
		}
		log.Errorf("database error, select by id: %v", err)
		return &model.Person{}, err
	}
	return &u, nil
}

// GetAllUsers select all users from db
func (p *PRepository) GetAllUsers(ctx context.Context) ([]*model.Person, error) {
	var persons []*model.Person
	rows, err := p.Pool.Query(ctx, "select id,name,works,age from persons")
	if err != nil {
		log.Errorf("database error with select all users, %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		per := model.Person{}
		err = rows.Scan(&per.ID, &per.Name, &per.Works, &per.Age)
		if err != nil {
			log.Errorf("database error with select all users, %v", err)
			return nil, err
		}
		persons = append(persons, &per)
	}

	return persons, nil
}

// DeleteUser delete user by id
func (p *PRepository) DeleteUser(ctx context.Context, id string) error {
	a, err := p.Pool.Exec(ctx, "delete from persons where id=$1", id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with this id doesnt exist: %v", err)
		}
		log.Errorf("error with delete user %v", err)
		return err
	}
	return nil
}

// UpdateUser update parameters for user
func (p *PRepository) UpdateUser(ctx context.Context, id string, per *model.Person) error {
	a, err := p.Pool.Exec(ctx, "update persons set name=$1,works=$2,age=$3 where id=$4", &per.Name, &per.Works, &per.Age, id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		log.Errorf("error with update user %v", err)
		return err
	}
	return nil
}

// UpdateAuth logout, delete refresh token
func (p *PRepository) UpdateAuth(ctx context.Context, id, refreshToken string) error {
	a, err := p.Pool.Exec(ctx, "update persons set refreshToken=$1 where id=$2", refreshToken, id)
	if a.RowsAffected() == 0 {
		return fmt.Errorf("user with this id doesnt exist")
	}
	if err != nil {
		log.Errorf("error with update user %v", err)
		return err
	}
	return nil
}

// SelectByIDAuth get id and refresh token by id
func (p *PRepository) SelectByIDAuth(ctx context.Context, id string) (model.Person, error) {
	per := model.Person{}
	err := p.Pool.QueryRow(ctx, "select id,refreshToken from persons where id=$1", id).Scan(&per.ID, &per.RefreshToken)

	if err != nil /*err==no-records*/ {
		if err == pgx.ErrNoRows {
			return model.Person{}, fmt.Errorf("user with this id doesnt exist: %v", err)
		}
		log.Errorf("database error, select by id: %v", err)
		return model.Person{}, err /*p, fmt.errorf("user with this id doesnt exist")*/
	}
	return per, nil
}
