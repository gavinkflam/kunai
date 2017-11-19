package core

import (
  "github.com/gavinkflam/kunai/configs"
)

func FullAssetPath(filename string) string {
  return configs.DirStr() + "/" + filename;
}
