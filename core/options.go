package core

import (
  "strconv"
)

type Options struct {
  // Size options
  Fit    string
  Width  int
  Height int

  // Format options
  Format string

  // Signature options
  Signature string
}

type fString func(string) string

func ParseOptions(p fString) (*Options, error) {
  o := new(Options)

  // Capture fit option and fallback to clip
  o.Fit = elemString(p, "fit", "clip")

  // Capture width option
  width, err := strconv.Atoi(elemString(p, "w", "0"))
  if err != nil { return nil, err }
  o.Width = width

  // Capture height option
  height, err := strconv.Atoi(elemString(p, "h", "0"))
  if err != nil { return nil, err }
  o.Height = height

  // Capture format option
  o.Format = p("fm")

  // Capture signature option
  o.Signature = p("s")

  return o, nil
}

func elemString(p fString, key, fallback string) string {
  if val := p(key); len(val) > 0 {
    return val
  } else {
    return fallback
  }
}
