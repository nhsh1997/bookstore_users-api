package users

import (
	"fmt"
	"github.com/nhsh1997/bookstore_users-api/datasources/mysql/users_db"
	"github.com/nhsh1997/bookstore_users-api/utils/errors"
	"github.com/nhsh1997/bookstore_users-api/utils/mysql_utils"
)

var (
	userDB = make(map[int64]*User)
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?);"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

func (user *User) Get() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()

	result := statement.QueryRow(user.ID)

	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}
func (user *User) Save() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()

	insertResult, saveErr := statement.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userID, err := insertResult.LastInsertId()

	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	user.ID = userID

	return nil
}

func (user *User) Update() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()

	_, updateErr := statement.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if updateErr != nil {
		return mysql_utils.ParseError(updateErr)
	}

	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to update user: %s", err.Error()))
	}

	return nil
}