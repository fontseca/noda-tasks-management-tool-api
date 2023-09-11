package main

import (
	"log"
	"noda/engine"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	engine.Run()
}
