package handler

import (
	"github.com/gin-gonic/gin"
)

func HTTPErrorHandler(c *gin.Context) {
	c.Next()

	if c.Errors.Last() != nil {
		c.JSON(-1, c.Errors.Last())
	}

}
