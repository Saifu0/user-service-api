package user

import (
	"database/sql"
	"fmt"
	"github.com/Saifu0/user-service-api/common/dates"
	usersdb "github.com/Saifu0/user-service-api/datasources/mysql/users_db"
	"strings"

	"github.com/Saifu0/user-service-api/common/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	erroNoRow        = "no rows in result set"
	queryInsertUser  = "INSERT INTO users (first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), erroNoRow) {
			return errors.NewNotFound(fmt.Sprintf("user with user id %d, not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while getting user with id %d: %s", user.Id, err.Error()))
	}

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
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequest(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to insert user: %s", err.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while getting last inserted ID: %s", err.Error()))
	}
	user.Id = userId
	return nil
}
