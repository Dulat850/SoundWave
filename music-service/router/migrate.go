package router

import (
	"errors"

	"music-service/config"
	"music-service/repositories"
)

func Migrate() error {
	cfg := config.Load()
	repositories.ConnectSQL(cfg)

	if repositories.SQLDB == nil {
		return errors.New("repositories.SQLDB is nil (ConnectSQL did not initialize db)")
	}

	return repositories.ApplyMigrations(repositories.SQLDB, "migrations")
}
