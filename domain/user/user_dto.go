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

func (user *User) Validate() *errors.RestErr {
	email := strings.TrimSpace(strings.ToLower(user.Email))
	if email == "" {
		return errors.NewBadRequest("invalid email address")
	}
	return nil
}
