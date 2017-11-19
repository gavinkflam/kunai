package apis

import (
  "github.com/gin-gonic/gin"
)

func internalErrorIfAny(c *gin.Context, err error) {
  if err != nil {
    c.AbortWithError(500, err)
  }
}
