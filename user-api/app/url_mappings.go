package app

import (
	log "github.com/sirupsen/logrus"

	userController "user-api/controller"
)

func mapUrls() {

	// Users Mapping
	router.GET("/user-api/user/:id", userController.GetUserById)
	router.GET("/user-api/user", userController.GetUsers)
	router.POST("/user-api/user", userController.UserInsert) // Sign In
	router.DELETE("user-api/user/:id", userController.DeleteUser)
	router.PUT("user-api/user/:id", userController.UpdateUser)

	log.Info("Finishing mappings configurations")
}
