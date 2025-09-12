package server

import (
	"blogs_learn/models"
	"blogs_learn/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type postBody struct {
	ID      uint   `json:"id"`
	Title   string `binding:"required" json:"title"`
	Content string `binding:"required" json:"content"`
}

func QueryPostList(c *gin.Context) {
	var postList []models.Post
	db := utils.GetDB(c)
	db.Select("id, title, user_id, created_at, updated_at").Find(&postList)
	c.JSON(http.StatusOK, postList)
}

func QueryPostDetail(c *gin.Context) {
	postID := c.Query("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	post := &models.Post{}
	db := utils.GetDB(c)
	db.Select("id, title, content, user_id, created_at, updated_at").First(post, postID)
	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	body := &postBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	userID := utils.GetLoginUserID(c)
	post := &models.Post{Title: body.Title, Content: body.Content, UserID: userID}
	db := utils.GetDB(c)
	db.Create(post)
	c.JSON(http.StatusOK, post)
}

func ModifyPost(c *gin.Context) {
	postBody := &postBody{}
	err := c.ShouldBindJSON(postBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	if postBody.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	db := utils.GetDB(c)
	result := db.Model(&models.Post{}).
		Where("id = ? and user_id = ?", postBody.ID, utils.GetLoginUserID(c)).
		Updates(map[string]any{
			"title":   postBody.Title,
			"content": postBody.Content,
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post updated"})
}

func DeletePost(c *gin.Context) {
	postID := c.Query("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}
	db := utils.GetDB(c)
	id, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	result := db.Model(&models.Post{}).
		Where("id = ? and user_id = ?", id, utils.GetLoginUserID(c)).
		Delete(&models.Post{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
