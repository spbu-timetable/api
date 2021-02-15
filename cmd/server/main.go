package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/spbu-timetable/api/internal/server"
)

func main() {
	godotenv.Load("config.env")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
