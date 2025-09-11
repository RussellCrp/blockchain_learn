package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet("db").(*gorm.DB)
}

func GetLoginUserID(c *gin.Context) uint {
	userID := c.MustGet("userID")
	return userID.(uint)
}

func GetToken(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}
