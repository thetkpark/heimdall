package handler

import "github.com/gin-gonic/gin"

func NewGinErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"message": message,
		"status":  code,
	})
}
