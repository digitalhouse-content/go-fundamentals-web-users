package user

import (
	"log"

	"context"

	"github.com/digitalhouse-content/go-fundamentals-web-users/internal/domain"
)

type DB struct {
	Users     []domain.User
	MaxUserID uint64
}

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
	}

	repo struct {
		db  DB
		log *log.Logger
	}
)

func NewRepo(db DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {

	r.db.MaxUserID++
	user.ID = r.db.MaxUserID
	r.db.Users = append(r.db.Users, *user)
	r.log.Println("repository create")
	return nil

}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("repository get all")
	return r.db.Users, nil
}
