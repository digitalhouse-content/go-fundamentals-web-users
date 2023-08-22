package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")

type ErrNotFound struct {
	ID uint64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user id '%d' doesn't exist", e.ID)
}
