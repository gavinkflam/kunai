package apis

import (
  "github.com/gin-gonic/gin"
)

func RegisterImagesApis(router *gin.RouterGroup) {
  router.GET("/*filename", ProcessImage)
}

func ProcessImage(c *gin.Context) {
  c.String(200, "It works!")
}
