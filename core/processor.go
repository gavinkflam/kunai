package core

import (
  "io/ioutil"
  "os"

  "github.com/h2non/bimg"
)

const tmpFilePrefix = "kunai"

func Process(filename string) (string, error) {
  fullFilename := FullAssetPath(filename)

  buffer, err := bimg.Read(fullFilename)
  if err != nil { return "", err }

  newImage, err := bimg.NewImage(buffer).Convert(bimg.PNG)
  if err != nil { return "", err }

  tmpFile, err := ioutil.TempFile(os.TempDir(), tmpFilePrefix)
  if err != nil { return "", err }

  err = bimg.Write(tmpFile.Name(), newImage)
  if err != nil { return "", err }

  return tmpFile.Name(), nil
}
