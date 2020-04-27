package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nhsh1997/bookstore_oauth-go/oauth"
	"github.com/nhsh1997/bookstore_users-api/domain/users"
	"github.com/nhsh1997/bookstore_users-api/services"
	"github.com/nhsh1997/bookstore_users-api/utils/errors"
	"net/http"
	"strconv"
)


func getUserId(userIdParam string)(int64, *errors.RestError){
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)

	if userErr != nil {
		// #TODO handle json error
		err := errors.NewBadRequestError("user id should be a number")
		return 0, err
	}
	return userId, nil
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}

	if callerId := oauth.GetCallerId(c.Request); callerId == 0 {
		err := errors.RestError{
			Message: "resource not available",
			Status:  http.StatusUnauthorized,
		}
		c.JSON(err.Status, err)
		return
	}
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	user, getErr := services.UsersServices.GetUser(userId)
	if getErr != nil {
		// #TODO handle get error
		c.JSON(getErr.Status, getErr)
		return
	}
	if oauth.GetCallerId(c.Request) == user.ID {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Create(c *gin.Context)  {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// #TODO handle json error
		restError := errors.NewBadRequestError("invalid json value")
		c.JSON(restError.Status, restError)
		return
	}

	result, saveErr := services.UsersServices.CreateUser(user)
	if saveErr != nil {
		// #TODO handle save error
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context)  {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)

	if userErr != nil {
		// #TODO handle json error
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// #TODO handle json error
		restError := errors.NewBadRequestError("invalid json value")
		c.JSON(restError.Status, restError)
		return
	}

	user.ID = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, err := services.UsersServices.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
	}
	_, getErr := services.UsersServices.GetUser(userId)
	if getErr != nil {
		// #TODO handle get error
		c.JSON(getErr.Status, getErr)
		return
	}

	if err := services.UsersServices.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status":"deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UsersServices.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context)  {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user, err := services.UsersServices.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}