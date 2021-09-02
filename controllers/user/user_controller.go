package user

import (
	"net/http"
	"strconv"

	"github.com/Saifu0/user-service-api/common/errors"
	"github.com/Saifu0/user-service-api/domain/user"
	"github.com/Saifu0/user-service-api/services"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		restErr := errors.NewBadRequest("invalid json body")
		c.JSON(http.StatusBadRequest, restErr)
		return
	}
	result, saveErr := services.CreateUser(u)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		userErr := errors.NewBadRequest("user id should be a number")
		c.JSON(http.StatusBadRequest, userErr)
		return
	}
	result, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func UpdateUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		userErr := errors.NewBadRequest("user id should be a number")
		c.JSON(http.StatusBadRequest, userErr)
		return
	}
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		restErr := errors.NewBadRequest("invalid json body")
		c.JSON(http.StatusBadRequest, restErr)
		return
	}

	isPartial := c.Request.Method == http.MethodPatch

	u.Id = userId
	result, updateErr := services.UpdateUser(isPartial, u)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, result)
}
