package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/env"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
)

func server(queries *db.Queries) error {
	e := echo.New()

	api := e.Group("/api/v1")
	api.GET("/photos", func(c echo.Context) error {
		data, err := queries.GetAllPhotosWithTags(c.Request().Context())
		if err != nil {
			slog.Error("get all photos with tags", "error", err)
			return c.String(500, "internal server error")
		}
		return c.JSON(200, data)
	})
	api.GET("/photos/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			slog.Error("parse photo id", "error", err)
			return c.String(400, "bad request: invalid photo id")
		}
		data, err := queries.GetPhotoByIDWithTags(c.Request().Context(), int32(id))
		if err != nil {
			slog.Error("get photo by id", "error", err)
			return c.String(500, "internal server error")
		}
		return c.JSON(200, data)
	})
	api.POST("/photos", handler(func(c echo.Context, req struct {
		Title    string         `json:"title" required:"false"`
		PhotoURL string         `json:"photo_url"`
		Comment  string         `json:"comment" required:"false"`
		Metadata map[string]any `json:"metadata" required:"false"`
	}) error {
		metadata, err := json.Marshal(req.Metadata)
		if err != nil {
			slog.Error("marshal metadata", "error", err)
			return c.String(400, "bad request: invalid metadata")
		}
		if err := queries.AddPhoto(context.Background(), db.AddPhotoParams{
			Title:    pgText(req.Title),
			PhotoUrl: req.PhotoURL,
			Comment:  pgText(req.Comment),
			Metadata: metadata,
		}); err != nil {
			slog.Error("add photo", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	}), RequireAdminMiddleware)

	api.GET("/tags", func(c echo.Context) error {
		data, err := queries.ListTagsWithCount(c.Request().Context())
		if err != nil {
			slog.Error("list tags with count", "error", err)
			return c.String(500, "internal server error")
		}
		return c.JSON(200, data)
	})
	api.POST("/tags", handler(func(c echo.Context, req struct {
		Title   string `json:"title"`
		Comment string `json:"comment" required:"false"`
	}) error {
		if err := queries.CreateTag(context.Background(), db.CreateTagParams{
			Title:   req.Title,
			Comment: pgText(req.Comment),
		}); err != nil {
			slog.Error("create tag", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	}))
	api.DELETE("/tags", handler(func(c echo.Context, req struct {
		Title string `json:"title"`
	}) error {
		if err := queries.DeleteTag(context.Background(), req.Title); err != nil {
			slog.Error("delete tag by title", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	}))

	api.GET("/admin", func(c echo.Context) error {
		return c.JSON(200, map[string]bool{
			"is_admin": isAdmin(c),
		})
	})
	api.POST("/admin", handler(func(c echo.Context, req struct {
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
	}), NewRateLimiter(5, 10*time.Minute, false).Middleware) // global rate limit 5 requests per 10 minutes for TOTP attempts

	return e.Start(env.DefaultEnv.ADDR)
}
