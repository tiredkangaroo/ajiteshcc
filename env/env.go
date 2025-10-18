package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	POSTGRES_CONNECTION_URI string
	ADDR                    string
}

var DefaultEnv Environment

func init() {
	godotenv.Load()
	DefaultEnv = Environment{
		POSTGRES_CONNECTION_URI: envRequire("POSTGRES_CONNECTION_URI"),
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
