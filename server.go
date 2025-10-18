package main

import (
	"context"
	"log/slog"

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
	api.POST("/photos", handler(func(c echo.Context, req struct {
		Title string `json:"title"`
	}) error {
		return nil
	}))

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
		return nil
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
	return e.Start(env.DefaultEnv.ADDR)
}
