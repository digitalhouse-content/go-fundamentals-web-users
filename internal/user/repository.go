package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/digitalhouse-content/go-fundamentals-web-users/internal/domain"
)

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
	}

	repo struct {
		db  *sql.DB
		log *log.Logger
	}
)

func NewRepo(db *sql.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {

	sqlQ := "INSERT INTO users(first_name, last_name, email) VALUES(?,?,?)"
	res, err := r.db.Exec(sqlQ, user.FirstName, user.LastName, user.Email)
	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	user.ID = uint64(id)
	r.log.Println("user created with id: ", id)

	return nil

}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	sqlQ := "SELECT id, first_name, last_name, email FROM users"

	rows, err := r.db.Query(sqlQ)
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {

			r.log.Println(err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	r.log.Println("user get all: ", len(users))
	return users, nil
}

func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	sqlQ := "SELECT id, first_name, last_name, email FROM users WHERE id = ?"
	var u domain.User

	if err := r.db.QueryRow(sqlQ, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
		r.log.Println(err.Error())
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{id}
		}
		return nil, err
	}

	r.log.Println("get user with id: ", id)
	return &u, nil
}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	var fields []string
	var values []interface{}

	if firstName != nil {
		fields = append(fields, "first_name=?")
		values = append(values, *firstName)
	}

	if lastName != nil {
		fields = append(fields, "last_name=?")
		values = append(values, *lastName)
	}

	if email != nil {
		fields = append(fields, "email=?")
		values = append(values, *email)
	}

	if len(fields) == 0 {
		r.log.Println(ErrThereArentFields.Error())
		return ErrThereArentFields
	}

	values = append(values, id)

	sqlQ := fmt.Sprintf("UPDATE users SET %s WHERE id=?", strings.Join(fields, ","))
	res, err := r.db.Exec(sqlQ, values...)
	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	if row == 0 {
		err := ErrNotFound{id}
		r.log.Println(err.Error())
		return err
	}

	r.log.Println("user updated id: ", id)
	return nil

}
