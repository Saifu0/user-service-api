package services

import (
	"github.com/Saifu0/user-service-api/common/errors"
	"github.com/Saifu0/user-service-api/domain/user"
)

func CreateUser(user user.User) (*user.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(userId int64) (*user.User, *errors.RestErr) {
	user := &user.User{Id: userId}
	if err := user.Get(); err != nil {
		return nil, err
	}
	return user, nil
}
