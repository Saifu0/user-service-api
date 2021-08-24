package app

import (
	"github.com/Saifu0/user-service-api/controllers/pong"
	"github.com/Saifu0/user-service-api/controllers/user"
)

func mapUrls() {
	router.GET("/ping", pong.Pong)

	router.GET("/users/:user_id", user.GetUser)
	router.POST("/user", user.CreateUser)
}
