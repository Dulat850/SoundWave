package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	MongoURI string

	StorageAudioDir  string
	StorageCoversDir string
	StorageTmpDir    string
	PublicBaseURL    string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found, using system env")
	}

	return &Config{
		Port: os.Getenv("PORT"),

		PostgresHost:     os.Getenv("DB_HOST"),
		PostgresPort:     os.Getenv("DB_PORT"),
		PostgresUser:     os.Getenv("DB_USER"),
		PostgresPassword: os.Getenv("DB_PASSWORD"),
		PostgresDB:       os.Getenv("DB_NAME"),

		MongoURI: os.Getenv("MONGODB_URI"),

		StorageAudioDir:  envOrDefault("STORAGE_AUDIO_DIR", "storage/audio"),
		StorageCoversDir: envOrDefault("STORAGE_COVERS_DIR", "storage/covers"),
		StorageTmpDir:    envOrDefault("STORAGE_TMP_DIR", "storage/tmp"),
		PublicBaseURL:    os.Getenv("PUBLIC_BASE_URL"),
	}
}

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
