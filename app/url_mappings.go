package app

import "github.com/nhsh1997/bookstore_users-api/controllers"

func mapUrls()  {
	router.GET("/users/:user_id", controllers.Get)
	router.GET("/users", controllers.Search)
	router.POST("/users", controllers.Create)
	router.PUT("/users/:user_id", controllers.Update)
	router.PATCH("/users/:user_id", controllers.Update)
	router.DELETE("/users/:user_id", controllers.Delete)
}
