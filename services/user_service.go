package services

import (
	"fmt"
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
	u := &user.User{Id: userId}
	if err := u.Get(); err != nil {
		return nil, err
	}
	return u, nil
}

func UpdateUser(isPartial bool, u user.User) (*user.User, *errors.RestErr) {
	current, err := GetUser(u.Id)
	if err != nil {
		return nil, err
	}
	fmt.Println(u)

	if isPartial {
		if u.FirstName != "" {
			current.FirstName = u.FirstName
		}
		if u.LastName != "" {
			current.LastName = u.LastName
		}
		if u.Email != "" {
			current.Email = u.Email
		}
	} else {
		current.FirstName = u.FirstName
		current.LastName = u.LastName
		current.Email = u.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}
