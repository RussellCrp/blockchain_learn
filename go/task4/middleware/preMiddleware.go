package middleware

import (
	"blogs_learn/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func PreMiddleware(c *gin.Context) {
	path := c.FullPath()
	m := make(map[string]any)
	params := c.Params
	fmt.Println("request path:", path)
	fmt.Println("request body:", m)
	fmt.Println("request param:", params)
	c.Set("db", config.DB)
	c.Next()
}
