package utils

import (
  "math/rand"
  "os"
  "time"
)

func GetEnv(key, fallback string) string {
  value := os.Getenv(key)

  if len(value) == 0 {
    return fallback
  }

  return value
}

func RandomInt(min int, max int) int {
  rand.Seed(time.Now().UTC().UnixNano())

  return min + rand.Intn(max-min)
}
