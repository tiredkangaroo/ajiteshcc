package server

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
)

func (s *Server) listTags(c echo.Context) error {
	data, err := s.Queries.ListTagsWithCount(c.Request().Context())
	if err != nil {
		slog.Error("list tags with count", "error", err)
		return c.String(500, "internal server error")
	}
	return c.JSON(200, data)
}

func (s *Server) addTagHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Title   string `json:"title"`
		Comment string `json:"comment" required:"false"`
	}) error {
		if err := s.Queries.CreateTag(context.Background(), db.CreateTagParams{
			Title:   req.Title,
			Comment: pgText(req.Comment),
		}); err != nil {
			slog.Error("create tag", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	})
}

func (s *Server) deleteTagHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Title string `json:"title"`
	}) error {
		if err := s.Queries.DeleteTag(context.Background(), req.Title); err != nil {
			slog.Error("delete tag by title", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	})
}
