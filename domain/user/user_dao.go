package user

import (
	"database/sql"
	"fmt"
	"github.com/Saifu0/user-service-api/common/dates"
	usersdb "github.com/Saifu0/user-service-api/datasources/mysql/users_db"

	"github.com/Saifu0/user-service-api/common/errors"
)

const (
	queryInsertUser = "INSERT INTO users (first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	if err := usersdb.Client.Ping(); err != nil {
		panic(err)
	}
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
	stmt, err := usersdb.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	user.DateCreated = dates.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to insert user: %s", err.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while getting last inserted ID: %s", err.Error()))
	}
	user.Id = userId
	return nil
}
