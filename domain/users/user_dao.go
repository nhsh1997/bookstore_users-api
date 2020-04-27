package users

import (
	"fmt"
	"github.com/nhsh1997/bookstore_users-api/datasources/mysql/users_db"
	"github.com/nhsh1997/bookstore_users-api/logger"
	"github.com/nhsh1997/bookstore_users-api/utils/errors"
	"github.com/nhsh1997/bookstore_users-api/utils/mysql_utils"
	"strings"
)

var (
	userDB = make(map[int64]*User)
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?"
	queryGetUser = "SELECT id, first_name, last_name, email, status, date_created FROM users WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id=?;"
	queryFindByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)


func (user *User) Get() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer statement.Close()

	result := statement.QueryRow(user.ID)

	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); err != nil {
		logger.Error("error when trying to prepare get user by id", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}
func (user *User) Save() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer statement.Close()

	insertResult, saveErr := statement.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to prepare save user", err)
		return errors.NewInternalServerError("database error")
	}

	userID, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}

	user.ID = userID

	return nil
}

func (user *User) Update() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer statement.Close()

	_, updateErr := statement.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if updateErr != nil {
		logger.Error("error when trying to prepare update user", err)
		return errors.NewInternalServerError("database error")
	}

	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to update user: %s", err.Error()), err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *errors.RestError  {
	statement, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer statement.Close()

	if _, deleteErr := statement.Exec(user.ID); err != nil {
		logger.Error("error when trying to prepare delete user", deleteErr)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestError) {
	statement, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer statement.Close()

	rows, err := statement.Query(status)
	if err != nil {
		logger.Error("error when trying to prepare find users by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next(){
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *errors.RestError {
	statement, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer statement.Close()

	result := statement.QueryRow(user.Email, user.Password, StatusActive)

	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), mysql_utils.ErrorNoRows){
			return errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to prepare get user by email and password", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}