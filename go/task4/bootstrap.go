package main

import (
	"blogs_learn/config"
	"blogs_learn/router"

	"github.com/gin-gonic/gin"
)

func init() {
	config.InitDB()
}

func main() {
	r := gin.Default()
	router.InitRouter(r)
	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
