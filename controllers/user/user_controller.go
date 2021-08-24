package user

import (
	"net/http"

	"github.com/Saifu0/user-service-api/common/errors"
	"github.com/Saifu0/user-service-api/domain/user"
	"github.com/Saifu0/user-service-api/services"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

func CreateUser(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequest("invalid json body")
		c.JSON(http.StatusBadRequest, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}
