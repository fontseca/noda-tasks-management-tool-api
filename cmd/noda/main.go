package main

import (
	"log"
	"noda/server"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	server.Run()
}
