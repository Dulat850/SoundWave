package main

import (
	"log"
	"music-service/config"
	"music-service/repositories"
	"music-service/router"

	"github.com/gin-gonic/gin"
)

func init() {
	cfg := config.Load()
	repositories.ConnectSQL(cfg)
}

func main() {

	r := gin.Default()

	if err := router.SetupRoutes(r); err != nil {
		log.Fatal(err)
	}

	_ = r.Run()
}
