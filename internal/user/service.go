package user

import (
	"context"
	"log"

	"github.com/digitalhouse-content/go-fundamentals-web-users/internal/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {

	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.log.Println("service create")
	return user, nil

}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	s.log.Println("service get all")
	return users, nil
}
