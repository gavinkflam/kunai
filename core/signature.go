package core

import (
  "crypto/md5"
  "encoding/hex"
  "regexp"
)

func CheckSignature(url, token string, option Options) bool {
  return option.Signature == deriveSignature(url, token)
}

func deriveSignature(url, token string) string {
  hasher := md5.New()
  hasher.Write([]byte(token + url))
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
