package services

import (
	"github.com/nhsh1997/bookstore_users-api/domain/users"
	"github.com/nhsh1997/bookstore_users-api/utils/crypto_utils"
	"github.com/nhsh1997/bookstore_users-api/utils/date_utils"
	"github.com/nhsh1997/bookstore_users-api/utils/errors"
)

var (
	UsersServices userServiceInterface = &userService{}
)

type userService struct {
}

type userServiceInterface interface {
	GetUser(int64)  (*users.User, *errors.RestError)
	CreateUser(users.User) (*users.User, *errors.RestError)
	UpdateUser(bool, users.User) (*users.User, *errors.RestError)
	DeleteUser(int64) *errors.RestError
	Search(string) (users.Users, *errors.RestError)
	LoginUser(users.LoginRequest)  (*users.User, *errors.RestError)
}


func (s *userService) GetUser(userId int64)  (*users.User, *errors.RestError)  {
	result := &users.User{
		ID: userId,
	}
	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}

func  (s *userService) CreateUser(user users.User) (*users.User, *errors.RestError)  {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.DateCreated = date_utils.GetNowDbFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func  (s *userService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestError)  {
	current, err := s.GetUser(user.ID)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func  (s *userService) DeleteUser(userId int64) *errors.RestError {
	user := &users.User{ID: userId}
	if err := user.Delete(); err != nil {
		return err
	}
	return nil
}

func  (s *userService) Search(status string) (users.Users, *errors.RestError)  {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func  (s *userService) LoginUser(request users.LoginRequest)  (*users.User, *errors.RestError)  {
	user := &users.User{
		Email: request.Email,
		Password: request.Password,
	}
	if err := user.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return user, nil
}