package server

import (
	"blogs_learn/models"
	"blogs_learn/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type commentBody struct {
	Content string `binding:"required"`
	PostID  uint   `binding:"required"`
}

func CreateComment(c *gin.Context) {
	body := &commentBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	userID := utils.GetLoginUserID(c)
	comment := &models.Comment{Content: body.Content, UserID: userID, PostID: body.PostID}
	db := utils.GetDB(c)
	db.Create(comment)
	c.JSON(http.StatusOK, comment)
}

func QueryCommentList(c *gin.Context) {
	postID := c.Query("postID")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	atoi, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is valid failed"})
	}
	commentList := []models.Comment{}
	db := utils.GetDB(c)
	db.Model(&models.Comment{}).Where("post_id = ? ", atoi).Find(&commentList)
	c.JSON(http.StatusOK, commentList)
}
