package env

import (
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	POSTGRES_CONNECTION_URI     string
	JWT_SECRET                  []byte
	R2_ACCOUNT_ID               string
	R2_ACCESS_KEY_ID            string
	R2_SECRET_ACCESS_KEY        string
	R2_PHOTOS_BUCKET_NAME       string
	R2_PHOTOS_BUCKET_PUBLIC_URL *url.URL
	TOTP_SECRET                 string
	DEBUG                       bool
	ADDR                        string
}

var DefaultEnv Environment

func init() {
	godotenv.Load()
	DefaultEnv = Environment{
		POSTGRES_CONNECTION_URI:     envRequire("POSTGRES_CONNECTION_URI"),
		JWT_SECRET:                  []byte(envRequire("JWT_SECRET")),
		TOTP_SECRET:                 envRequire("TOTP_SECRET"),
		R2_ACCOUNT_ID:               envRequire("R2_ACCOUNT_ID"),
		R2_ACCESS_KEY_ID:            envRequire("R2_ACCESS_KEY_ID"),
		R2_SECRET_ACCESS_KEY:        envRequire("R2_SECRET_ACCESS_KEY"),
		R2_PHOTOS_BUCKET_NAME:       envDefault("R2_PHOTOS_BUCKET_NAME", "photos"),
		R2_PHOTOS_BUCKET_PUBLIC_URL: urlRequire(envRequire("R2_PHOTOS_BUCKET_PUBLIC_URL")),
		DEBUG:                       os.Getenv("DEBUG") == "true",
		ADDR:                        envRequire("ADDR"),
	}
}

func urlRequire(value string) *url.URL {
	v, err := url.Parse(value)
	if err != nil {
		panic("invalid URL: " + value)
	}
	return v
}

func envRequire(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("environment variable " + key + " is required but not set")
	}
	return value
}

func envDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
