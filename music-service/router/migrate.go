package router

import (
	"errors"

	"music-service/config"
	"music-service/models"
	"music-service/repositories"
)

func Migrate() error {
	_ = config.Load()
	repositories.ConnectDB()

	if repositories.DB == nil {
		return errors.New("repositories.DB is nil (ConnectDB did not initialize DB)")
	}

	return repositories.DB.AutoMigrate(&models.User{})
}
