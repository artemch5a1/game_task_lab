package config

import (
	"os"
	"strings"
)

type Config struct {
	DBPath     string
	JWTSecret  string
	JWTIssuer  string
	JWTTTLHours int
}

const defaultDBPath = "data/app.db"

// NOTE: dev default secret. Override via JWT_SECRET in production.
const defaultJWTSecret = "game_task_lab_dev_secret_2026_change_me"

func Load() Config {
	dbPath := strings.TrimSpace(os.Getenv("DB_PATH"))
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		secret = defaultJWTSecret
	}

	issuer := strings.TrimSpace(os.Getenv("JWT_ISSUER"))
	if issuer == "" {
		issuer = "game-task-lab"
	}

	ttlHours := 24

	return Config{
		DBPath:     dbPath,
		JWTSecret:  secret,
		JWTIssuer:  issuer,
		JWTTTLHours: ttlHours,
	}
}

