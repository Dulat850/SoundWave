package main

import (
	"log"

	"music-service/router"
)

func main() {
	if err := router.Migrate(); err != nil {
		log.Fatal(err)
	}
	log.Println("Migration completed")
}
