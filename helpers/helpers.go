package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromCookie(c *gin.Context) (int, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return -1, errors.New("userId not found in context")
	}

	userIdInt, ok := userId.(int)
	if !ok {
		return -1, errors.New("userId is not of type int")
	}
	return userIdInt, nil
}
