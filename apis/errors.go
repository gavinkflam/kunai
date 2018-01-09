package apis

import (
	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	c.AbortWithError(500, err)
	return true
}
