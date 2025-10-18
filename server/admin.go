package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/tiredkangaroo/ajiteshcc/env"
)

// errors
var InvalidTOTPError = fmt.Errorf("invalid TOTP code")

// TOTP validation options
var totp_validation_opts = totp.ValidateOpts{
	Period:    30,
	Skew:      1,
	Digits:    otp.DigitsEight,
	Algorithm: otp.AlgorithmSHA256,
	Encoder:   otp.EncoderDefault, // b32 i think
}

func (s *Server) isAdmin(c echo.Context) error {
	return c.JSON(200, map[string]bool{
		"is_admin": isAdmin(c),
	})
}
func (s *Server) adminLoginHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		TOTP string `json:"totp"`
	}) error {
		jwtToken, err := issueJWTWithTOTP(req.TOTP)
		if err != nil {
			slog.Error("issue jwt with totp", "error", err)
			return c.String(403, "forbidden: invalid TOTP code")
		}
		c.SetCookie(&http.Cookie{
			Name:     "admin_token",
			Value:    jwtToken,
			Expires:  time.Now().Add(1 * time.Hour),
			HttpOnly: true,
			Secure:   !env.DefaultEnv.DEBUG,
			SameSite: http.SameSiteLaxMode,
		})
		return c.NoContent(204)
	})
}

// middleware to require admin access
func RequireAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isAdmin(c) {
			return c.String(403, "forbidden: admin access required")
		}
		return next(c)
	}
}

// verify if the request has a valid admin JWT token
func isAdmin(c echo.Context) bool {
	adminToken, err := c.Cookie("admin_token")
	if err != nil {
		slog.Error("get admin token cookie", "error", err)
		return false
	}
	token, err := jwt.ParseWithClaims(adminToken.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(env.DefaultEnv.JWT_SECRET), nil
	})
	return err == nil && token.Valid
}

// issue a JWT token if the provided TOTP code is valid
func issueJWTWithTOTP(timePasscode string) (string, error) {
	now := time.Now()
	ok, err := totp.ValidateCustom(timePasscode, env.DefaultEnv.TOTP_SECRET, now, totp_validation_opts)
	if err != nil {
		return "", fmt.Errorf("validate totp: %w", err)
	} else if !ok {
		return "", InvalidTOTPError
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "ajiteshcc-admin-service",
		Subject:   "administrator",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(env.DefaultEnv.JWT_SECRET)
	return signedToken, err
}
