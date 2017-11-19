package configs

import (
  "os"
  "strconv"
)

func Port() int {
  port, _ := strconv.Atoi(PortStr())
  return port
}

func PortStr() string {
  return getenv("PORT", "8080")
}

func getenv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}
