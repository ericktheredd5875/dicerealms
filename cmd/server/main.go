package main

import (
	"log"

	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/ericktheredd5875/dicerealms/internal/server"
	"github.com/ericktheredd5875/dicerealms/pkg/utils"
)

func main() {
	log.Println("Starting DiceRealms server...")

	if err := db.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	port := utils.ObtainEnv("TELNET_PORT", "4000")
	err := server.Start(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
