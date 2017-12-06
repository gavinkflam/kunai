package core

import (
  "crypto/md5"
  "encoding/hex"
  "errors"
  "regexp"
)

func CheckSignature(host, path, token string, option *Options) error {
  if len(option.Signature) == 0 {
    return errors.New("signature is required")
  }
  if option.Signature != deriveSignature(host + path, token) {
    return errors.New("signature does not match") }

  return nil
}

func deriveSignature(url, token string) string {
  shreddedUrl := shredSignatureFromUrl(url)
  hasher := md5.New()
  hasher.Write([]byte(token + shreddedUrl))
  return hex.EncodeToString(hasher.Sum(nil))
}

var signaturePatternQu = regexp.MustCompile("\\?s=[a-z0-9]{32}$")
var signaturePatternAm = regexp.MustCompile("&s=[a-z0-9]{32}$")

func shredSignatureFromUrl(url string) string {
  if signaturePatternQu.MatchString(url) {
    return signaturePatternQu.ReplaceAllString(url, "")
  } else if signaturePatternAm.MatchString(url) {
    return signaturePatternAm.ReplaceAllString(url, "")
  }
  return url
}
