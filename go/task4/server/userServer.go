package server

import (
	"blogs_learn/models"
	"blogs_learn/security"
	"blogs_learn/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required" `
	Email    string `json:"email"`
}

func Login(c *gin.Context) {
	body := &UserBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	db := utils.GetDB(c)
	user, err := getUserByUserName(db, body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or Password"})
		return
	}

	token, err := security.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token generation failed", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
	body := &UserBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	db := utils.GetDB(c)
	_, err := getUserByUserName(db, body.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exist"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash Password"})
		return
	}
	user := &models.User{Username: body.Username, Password: string(hashedPassword), Email: body.Email}

	db.Create(user)
	c.JSON(http.StatusOK, user)
}

func getUserByUserName(db *gorm.DB, Username string) (*models.User, error) {
	user := &models.User{}
	err := db.Where("username = ?", Username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
