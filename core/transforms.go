package core

import (
  "errors"
  "fmt"

  "github.com/h2non/bimg"
)

func transformImage(image *bimg.Image, options *Options) ([]byte, error) {
  // Delegate to fit transform
  _, err := fitTransform(image, options)
  if err != nil { return nil, err }

  return image.Image(), err
}

func fitTransform(image *bimg.Image, options *Options) ([]byte, error) {
  if options.Fit == "clip" {
    return clipTransform(image, options)
  } else if options.Fit == "crop" {
    return cropTransform(image, options)
  }

  return nil, fmt.Errorf("fit %s not supported", options.Fit)
}

func clipTransform(image *bimg.Image, options *Options) ([]byte, error) {
  var err error
  w, h := options.Width, options.Height

  if options.Width > 0 && options.Height > 0 {
    // Calculate size to fit within the bound without cropping or distorting
    w, h, err = calcNewSize(image, w, h)
  } else if options.Width > 0 {
    // Calculate auto height to matain aspeect ratio
    h, err = calcNewHeight(image, options.Width)
  } else if options.Height > 0 {
    // Calculate auto width to matain aspeect ratio
    w, err = calcNewWidth(image, options.Height)
  } else {
    // Do nothing if no height or width were provided
    return image.Image(), nil
  }

  if err != nil { return nil, err }
  return image.Resize(w, h)
}

func cropTransform(image *bimg.Image, options *Options) ([]byte, error) {
  if options.Width > 0 && options.Height > 0 {
    return image.ResizeAndCrop(options.Width, options.Height)
  }
  return nil, errors.New("both width and height are required for crop")
}

func calcNewSize(image *bimg.Image, width, height int) (int, int, error) {
  ratio, err := aspectRatio(image)
  if err != nil { return 0, 0, err }

  // Calculate requested ratio to compare with image aspect ratio
  w, h     := width, height
  reqRatio := float64(width) / float64(height)

  if reqRatio < ratio {
    // Calculate new height as provided height is too large
    h = int(ratio / float64(width))
  } else if reqRatio > ratio {
    // Calculate new height as provided width is too large
    w = int(ratio * float64(height))
  }

  return w, h, nil
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
