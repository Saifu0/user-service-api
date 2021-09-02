package app

import (
	"github.com/Saifu0/user-service-api/controllers/pong"
	"github.com/Saifu0/user-service-api/controllers/user"
)

func mapUrls() {
	router.GET("/ping", pong.Pong)

	router.POST("/user", user.Create)
	router.GET("/users/:user_id", user.Get)
	router.PUT("/users/:user_id", user.Update)
	router.PATCH("/users/:user_id", user.Update)
	router.DELETE("/users/:user_id", user.Delete)
	router.GET("internal/users/search", user.Search)
}
