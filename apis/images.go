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
  internalErrorIfAny(c, err)

  if configs.SignatureRequired() {
    err := core.CheckSignature(
      configs.HostStr(), c.Request.RequestURI, configs.Token(), options)
    internalErrorIfAny(c, err)
  }

  tmpFileName, err := core.Process(filename, options)
  internalErrorIfAny(c, err)

  c.File(tmpFileName)
  os.Remove(tmpFileName)
}
