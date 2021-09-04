package user

import (
	"database/sql"
	"fmt"
	"strings"

	usersdb "github.com/Saifu0/user-service-api/datasources/mysql/users_db"
	"github.com/go-sql-driver/mysql"

	"github.com/Saifu0/user-service-api/common/errors"
)

const (
	erroNoRow             = "no rows in result set"
	queryInsertUser       = "INSERT INTO users (first_name, last_name, email, password, status, date_created) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, status, date_created FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
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
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); getErr != nil {
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

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.Status, user.DateCreated)
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

func (user *User) FindByUser(status string) ([]User, *errors.RestErr) {
	stmt, err := usersdb.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		fmt.Println("error!")
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error while trying to find users using status, %s", err.Error()))
	}

	result := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
			return nil, errors.NewInternalServerError(fmt.Sprintf("error while trying to find users using status, %s", err.Error()))
		}
		result = append(result, u)
	}
	if len(result) == 0 {
		return nil, errors.NewNotFound(fmt.Sprintf("no users with status %s found", status))
	}
	return result, nil
}
