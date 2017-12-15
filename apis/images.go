package apis

import (
  "os"

  "github.com/gavinkflam/kunai/core"
  "github.com/gavinkflam/kunai/configs"
  "github.com/gin-gonic/gin"
)

func RegisterImagesApis(router *gin.RouterGroup) {
  router.GET("/*filename", ProcessImage)
}

func ProcessImage(c *gin.Context) {
  filename := c.Param("filename")

  options, err := core.ParseOptions(c.Query)
  if abortWithError(c, err) { return }

  if configs.SignatureRequired() {
    err := core.CheckSignature(
      c.Request.RequestURI, configs.Token(), options)
    if abortWithError(c, err) { return }
  }

  tmpFileName, err := core.Process(filename, options)
  if abortWithError(c, err) { return }

  c.File(tmpFileName)
  os.Remove(tmpFileName)
}
