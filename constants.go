package noda

import (
	"log"
	"os"
	"strings"
	"time"
)

const JWTExpiresIn = 1 * time.Hour

const Blankset = "  \a\b\f\r\t\v  "

var secret string

func init() {
	secret = strings.Trim(os.Getenv("JWT_SECRET"), Blankset)
	if "" == secret {
		log.Fatal("could not load environment variable \"JWT_SECRET\"")
	}
}

func Secret() []byte {
	return []byte(secret)
}
