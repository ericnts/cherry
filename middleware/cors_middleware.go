package middleware

import (
	"github.com/gin-gonic/gin"
)

// 开启跨域函数
func CorsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE,HEAD")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
