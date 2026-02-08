package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"

	"music-service/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var SQLDB *sql.DB

func ConnectSQL(cfg *config.Config) {
	dsn, err := postgresDSN(cfg)
	if err != nil {
		log.Fatal("invalid postgres config:", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("failed to open SQL db:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping SQL db:", err)
	}

	SQLDB = db
}

func postgresDSN(cfg *config.Config) (string, error) {
	if cfg == nil {
		return "", errors.New("config is nil")
	}
	if cfg.PostgresHost == "" || cfg.PostgresPort == "" || cfg.PostgresUser == "" || cfg.PostgresDB == "" {
		return "", fmt.Errorf("missing postgres settings: host=%q port=%q user=%q db=%q",
			cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresDB,
		)
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.PostgresUser, cfg.PostgresPassword),
		Host:   fmt.Sprintf("%s:%s", cfg.PostgresHost, cfg.PostgresPort),
		Path:   cfg.PostgresDB,
	}

	q := url.Values{}
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	return u.String(), nil
}
