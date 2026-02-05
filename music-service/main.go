package main

import (
	"music-service/config"
	"music-service/repositories"
	"music-service/router"

	"github.com/gin-gonic/gin"
)

func init() {
	_ = config.Load()
	repositories.ConnectDB()
}

func main() {

	r := gin.Default()

	r.POST("/auth/signup", router.CreateUser)
	r.POST("/auth/login", router.Login)
	r.GET("/user/profile", router.CheckAuth, router.GetUserProfile)

	_ = r.Run()
}
