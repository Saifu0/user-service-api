package services

import (
	"github.com/Saifu0/user-service-api/common/errors"
	"github.com/Saifu0/user-service-api/domain/user"
)

func CreateUser(user user.User) (*user.User, *errors.RestErr) {
	return &user, nil
}
