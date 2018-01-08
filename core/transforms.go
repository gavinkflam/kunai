package core

import (
  "errors"
  "fmt"

  "github.com/h2non/bimg"
)

func transformImage(image *bimg.Image, options *Options) ([]byte, error) {
  var err error

  // Delegate to fit transform
  _, err = fitTransforms(image, options)
  if err != nil { return nil, err }

  // Delegate to format transforms
  _, err = formatTransforms(image, options)
  if err != nil { return nil, err }

  return image.Image(), nil
}

func fitTransforms(image *bimg.Image, options *Options) ([]byte, error) {
  if options.Fit == "clip" {
    return clipTransform(image, options)
  } else if options.Fit == "crop" {
    return cropTransform(image, options)
  }

  return nil, fmt.Errorf("fit %s not supported", options.Fit)
}

func formatTransforms(image *bimg.Image, options *Options) ([]byte, error) {
  var err error
  // Apply color space transform
  _, err = colorSpaceTransform(image, options)
  if err != nil { return nil, err }

  // Apply output format transform
  _, err = outputFormatTransform(image, options)
  if err != nil { return nil, err }

  // Apply quality format transform
  _, err = qualityTransform(image, options)
  if err != nil { return nil, err }

  return image.Image(), nil
}

func qualityTransform(image *bimg.Image, options *Options) ([]byte, error) {
  imageOptions := bimg.Options{
    Quality: options.Quality,
  }

  return image.Process(imageOptions)
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

func colorSpaceTransform(image *bimg.Image, options *Options) ([]byte, error) {
  if options.ColorSpace == "srgb" {
    return image.Colourspace(bimg.InterpretationSRGB)
  }
  if options.ColorSpace == "strip" {
    processOptions := bimg.Options{
      NoProfile: true,
    }
    return image.Process(processOptions)
  }
  return nil, fmt.Errorf("color space %s not supported", options.ColorSpace)
}

func outputFormatTransform(image *bimg.Image, options *Options) ([]byte, error) {
  if len(options.Format) == 0 {
    // Do nothing if output format was not supplied
    return image.Image(), nil
  }

  // Convert format to image type and apply
  imageType, err := formatToImageType(options.Format)
  if err != nil { return nil, err }

  return image.Convert(imageType)
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

func formatToImageType(fm string) (bimg.ImageType, error) {
  switch fm {
  case "png":
    return bimg.PNG, nil
  case "jpg":
    return bimg.JPEG, nil
  case "webp":
    return bimg.WEBP, nil
  case "gif":
    return bimg.GIF, nil
  default:
    return bimg.UNKNOWN, fmt.Errorf("format %s not supported", fm)
  }
}
