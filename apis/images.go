package apis

import (
  "os"

  "github.com/gavinkflam/kunai/core"
  "github.com/gin-gonic/gin"
)

func RegisterImagesApis(router *gin.RouterGroup) {
  router.GET("/*filename", ProcessImage)
}

func ProcessImage(c *gin.Context) {
  filename := c.Param("filename")
  tmpFileName, err := core.Process(filename)

  if err != nil {
    c.String(500, err.Error())
  }

  c.File(tmpFileName)
  os.Remove(tmpFileName)
}
