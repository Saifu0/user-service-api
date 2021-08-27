package user

import (
	"fmt"

	"github.com/Saifu0/user-service-api/common/errors"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	result := userDB[user.Id]
	if result == nil {
		return errors.NewNotFound(fmt.Sprintf("user %d not found", user.Id))
	}
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestErr {
	current := userDB[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequest(fmt.Sprintf("email %s already exists", current.Email))
		}
		return errors.NewBadRequest(fmt.Sprintf("user id %d already exists", user.Id))
	}
	userDB[user.Id] = user
	return nil
}
