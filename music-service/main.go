package main

import (
	"log"
	"music-service/config"
	"music-service/repositories"
	"music-service/router"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	cfg := config.Load()
	repositories.ConnectSQL(cfg)
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // адрес фронта
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if err := router.SetupRoutes(r); err != nil {
		log.Fatal(err)
	}

	_ = r.Run() // по умолчанию на :8080
}
