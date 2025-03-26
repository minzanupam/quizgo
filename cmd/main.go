package main

import (
	"log"
	"quizgo/src/services"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	services.HttpService()
}
