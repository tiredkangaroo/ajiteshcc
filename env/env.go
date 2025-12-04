package env

import (
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	// POSTGRES_CONNECTION_URI is the connection URI for the PostgreSQL database
	POSTGRES_CONNECTION_URI string
	// JWT_SECRET is the secret key used for signing JSON Web Tokens (JWTs)
	JWT_SECRET []byte
	// R2_ACCOUNT_ID is the account ID for R2 access (found in R2 dashboard)
	R2_ACCOUNT_ID string
	// R2_ACCESS_KEY_ID is the access key ID for R2 access (found when creating R2 API keys)
	R2_ACCESS_KEY_ID string
	// R2_SECRET_ACCESS_KEY is the secret access key for R2 access (in conjunction with R2_ACCESS_KEY_ID)
	R2_SECRET_ACCESS_KEY string
	// R2_PHOTOS_BUCKET_NAME is the name of the photos bucket in R2 (e.g., "photos")
	R2_PHOTOS_BUCKET_NAME string
	// R2_PHOTOS_BUCKET_PUBLIC_URL is the public URL of the photos bucket (e.g., https://photos.ajitesh.cc)
	R2_PHOTOS_BUCKET_PUBLIC_URL *url.URL
	// CERT_FILE is the path to the TLS certificate file
	CERT_FILE string
	// KEY_FILE is the path to the TLS key file
	KEY_FILE string
	// TOTP_SECRET is the totp secret used for admin login
	TOTP_SECRET string
	// DEBUG allows for insecure behaviors. DO NOT ENABLE IN PRODUCTION
	DEBUG bool
	// CORS_ALLOWED_ORIGINS is a comma-separated list of allowed origins for CORS requests.
	// this field is ignored if DEBUG is true, in which case all origins are allowed.
	CORS_ALLOWED_ORIGINS string
	// ADDR is the address the server listens on (e.g., ":8080")
	ADDR string
	// BACKEND_DOMAIN is the domain of the backend server (used for setting cookie domain)
	BACKEND_DOMAIN string
	// SPOTIFY_CLIENT_ID is the client ID for Spotify API
	SPOTIFY_CLIENT_ID string
	// SPOTIFY_CLIENT_SECRET is the client secret for Spotify API
	SPOTIFY_CLIENT_SECRET string
	// SPOTIFY_REDIRECT_URI is the redirect URI for Spotify OAuth
	SPOTIFY_REDIRECT_URI string
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
		CERT_FILE:                   envDefault("CERT_FILE", ""),
		KEY_FILE:                    envDefault("KEY_FILE", ""),
		R2_PHOTOS_BUCKET_PUBLIC_URL: urlRequire(envRequire("R2_PHOTOS_BUCKET_PUBLIC_URL")),
		DEBUG:                       os.Getenv("DEBUG") == "true",
		CORS_ALLOWED_ORIGINS:        envDefault("CORS_ALLOWED_ORIGINS", "https://ajitesh.cc"),
		ADDR:                        envRequire("ADDR"),
		BACKEND_DOMAIN:              envRequire("BACKEND_DOMAIN"),
		SPOTIFY_CLIENT_ID:           envRequire("SPOTIFY_CLIENT_ID"),
		SPOTIFY_CLIENT_SECRET:       envRequire("SPOTIFY_CLIENT_SECRET"),
		SPOTIFY_REDIRECT_URI:        envRequire("SPOTIFY_REDIRECT_URI"),
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
