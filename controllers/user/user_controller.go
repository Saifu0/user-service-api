package user

import (
	"net/http"
	"strconv"

	"github.com/Saifu0/user-service-api/common/errors"
	"github.com/Saifu0/user-service-api/domain/user"
	"github.com/Saifu0/user-service-api/services"
	"github.com/gin-gonic/gin"
)

func getUserId(userId string) (int64, *errors.RestErr) {
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		idErr := errors.NewBadRequest("user id should be a number")
		return 0, idErr
	}
	return id, nil
}

func Create(c *gin.Context) {
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
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, idErr)
		return
	}
	result, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, idErr)
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
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, idErr)
		return
	}
	if err := services.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users)
}
