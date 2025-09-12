package router

import (
	"blogs_learn/middleware"
	"blogs_learn/server"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.PreMiddleware)
	r.POST("/api/login", server.Login)
	r.POST("/api/register", server.Register)
	r.GET("/api/post/list", server.QueryPostList)
	r.GET("/api/post/detail", server.QueryPostDetail)
	r.GET("/api/comment/list", server.QueryCommentList)

	authGroup := r.Group("/api/auth")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.POST("/post/create", server.CreatePost)
	authGroup.PUT("/post/modify", server.ModifyPost)
	authGroup.DELETE("/post/delete", server.DeletePost)
	authGroup.POST("/comment/create", server.CreateComment)
}
