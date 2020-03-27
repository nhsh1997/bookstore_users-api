package users

import (
	"github.com/nhsh1997/bookstore_users-api/utils/errors"
	"strings"
)

type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
}

func (user *User) Validate() *errors.RestError {
	email := strings.TrimSpace(strings.ToLower(user.Email))
	if email == "" {
		return errors.NewBadRequestError("Email is invalid")
	}
	return nil
}