package main

import (
  "log"
  "os"

  "github.com/gavinkflam/kunai/apis"
  "github.com/gavinkflam/kunai/configs"
  "github.com/gin-gonic/gin"
)

func main() {
  // Use stdout to output logs
  log.SetOutput(os.Stdout)

  // Default gin router
  router := gin.Default()

  // Endpoint for images processing
  apis.RegisterImagesApis(router.Group("/"))

  // Use logger middleware
  router.Use(gin.Logger())

  // Run on configured port:
  router.Run(":" + configs.PortStr())
}
