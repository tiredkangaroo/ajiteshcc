package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/tiredkangaroo/ajiteshcc/env"
)

var totp_validation_opts = totp.ValidateOpts{
	Period:    30,
	Skew:      1,
	Digits:    otp.DigitsEight,
	Algorithm: otp.AlgorithmSHA256,
	Encoder:   otp.EncoderDefault, // b32 i think
}
var InvalidTOTPError = fmt.Errorf("invalid TOTP code")

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

func RequireAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isAdmin(c) {
			return c.String(403, "forbidden: admin access required")
		}
		return next(c)
	}
}
