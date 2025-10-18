package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	POSTGRES_CONNECTION_URI string
	JWT_SECRET              []byte
	TOTP_SECRET             string
	DEBUG                   bool
	ADDR                    string
}

var DefaultEnv Environment

func init() {
	godotenv.Load()
	DefaultEnv = Environment{
		POSTGRES_CONNECTION_URI: envRequire("POSTGRES_CONNECTION_URI"),
		JWT_SECRET:              []byte(envRequire("JWT_SECRET")),
		TOTP_SECRET:             envRequire("TOTP_SECRET"),
		DEBUG:                   os.Getenv("DEBUG") == "true",
		ADDR:                    envRequire("ADDR"),
	}
}

func envRequire(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("environment variable " + key + " is required but not set")
	}
	return value
}
