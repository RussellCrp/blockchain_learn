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
	r.POST("/api/comment/list", server.CreateComment)

	authGroup := r.Group("/api/auth")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.POST("/post/create", server.CreatePost)
	authGroup.POST("/post/modify", server.ModifyPost)
	authGroup.POST("/post/delete", server.DeletePost)
	authGroup.POST("/comment/create", server.CreateComment)
}
