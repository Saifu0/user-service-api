package user

import (
	"strings"

	"github.com/Saifu0/user-service-api/common/errors"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"-"`
}

type Users []User

func (user *User) Validate() *errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequest("invalid email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" || len(user.Password) < 8 {
		return errors.NewBadRequest("invalid password")
	}
	return nil
}
