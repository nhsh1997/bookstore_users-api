package app

import "github.com/nhsh1997/bookstore_users-api/controllers"

func mapUrls()  {
	router.GET("/users/:user_id", controllers.GetUser)
	router.GET("/users", controllers.SearchUser)
	router.POST("/users", controllers.CreateUser)
}
