package user

import (
	"database/sql"
	"fmt"
	"github.com/Saifu0/user-service-api/common/dates"
	usersdb "github.com/Saifu0/user-service-api/datasources/mysql/users_db"
	"github.com/go-sql-driver/mysql"
	"strings"

	"github.com/Saifu0/user-service-api/common/errors"
)

const (
	erroNoRow       = "no rows in result set"
	queryInsertUser = "INSERT INTO users (first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id=?;"
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
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); getErr != nil {
		if strings.Contains(getErr.Error(), erroNoRow) {
			return errors.NewNotFound(fmt.Sprintf("user with user id %d, not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while getting user with id %d: %s", user.Id, getErr.Error()))
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
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		sqlErr, ok := saveErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error while trying to insert user: %s", saveErr.Error()))
		}
		switch sqlErr.Number {
		case 1092:
			return errors.NewBadRequest(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to insert user: %s", saveErr.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while getting last inserted ID: %s", err.Error()))
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryUpdateUser)
	if err != nil {
		fmt.Println("error!")
		return errors.NewInternalServerError(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		sqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error while trying to update user: %s", err.Error()))
		}
		switch sqlErr.Number {
		case 1092:
			return errors.NewBadRequest(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to update user: %s", err.Error()))
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryDeleteUser)
	if err != nil {
		fmt.Println("error!")
		return errors.NewInternalServerError(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)
	_, err = stmt.Exec(user.Id)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error while trying to delete user: %s", err.Error()))
	}
	return nil
}
