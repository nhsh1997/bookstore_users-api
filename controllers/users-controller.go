package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nhsh1997/bookstore_users-api/domain/users"
	"github.com/nhsh1997/bookstore_users-api/services"
	"github.com/nhsh1997/bookstore_users-api/utils/errors"
	"net/http"
)

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement me!")
}

func SearchUser(c *gin.Context)  {
	c.String(http.StatusNotImplemented, "Implement me!")
}

func CreateUser(c *gin.Context)  {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// #TODO handle json error
		restError := errors.NewBadRequestError("invalid json value")
		c.JSON(restError.Status, restError)
		return
	}

	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		// #TODO handle save error
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}