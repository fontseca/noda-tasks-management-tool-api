package global

import (
  "log"
  "os"
  "strings"
)

var secret string

func init() {
	secret = strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if "" == secret {
		log.Fatal("could not load env var: JWT_SECRET")
	}
}

func Secret() []byte {
	return []byte(secret)
}
