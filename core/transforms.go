package core

import (
  "errors"
  "fmt"

  "github.com/h2non/bimg"
)

func transformImage(image *bimg.Image, options *Options) ([]byte, error) {
  var rImage []byte
  var err    error

  // Delegate to fit transforms
  if options.Fit == "clip" {
    rImage, err = clipTransform(image, options)
  } else if options.Fit == "crop" {
    rImage, err = cropTransform(image, options)
  } else {
    return nil, fmt.Errorf("fit %s not supported", options.Fit)
  }

  return rImage, err
}

func clipTransform(image *bimg.Image, options *Options) ([]byte, error) {
  if options.Width > 0 && options.Height > 0 {
    return image.Resize(options.Width, options.Height)
  } else if options.Width > 0 {
    height, err := calcNewHeight(image, options.Width)

    if err != nil { return nil, err }
    return image.Resize(options.Width, height)
  } else {
    width, err := calcNewWidth(image, options.Height)

    if err != nil { return nil, err }
    return image.Resize(width, options.Height)
  }
}

func cropTransform(image *bimg.Image, options *Options) ([]byte, error) {
  if options.Width > 0 && options.Height > 0 {
    return image.ResizeAndCrop(options.Width, options.Height)
  } else {
    return nil, errors.New("both width and height are required for crop")
  }
}

func calcNewWidth(image *bimg.Image, height int) (int, error) {
  ratio, err := aspectRatio(image)
  if err != nil { return 0, err }

  return int(ratio * float64(height)), nil
}

func calcNewHeight(image *bimg.Image, width int) (int, error) {
  ratio, err := aspectRatio(image)
  if err != nil { return 0, err }

  return int(ratio / float64(width)), nil
}

func aspectRatio(image *bimg.Image) (float64, error) {
  size, err := image.Size()
  if err != nil { return 0, err }

  return float64(size.Width) / float64(size.Height), nil
}
