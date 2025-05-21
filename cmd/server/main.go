package main

import (
	"log"

	"github.com/ericktheredd5875/dicerealms/internal/server"
)

func main() {
	log.Println("Starting DiceRealms server...")
	err := server.Start(":4000")
	if err != nil {
		log.Fatal(err)
	}
}
