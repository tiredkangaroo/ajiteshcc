package server

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
)

// GET /api/v1/tags
func (s *Server) listTags(c echo.Context) error {
	data, err := s.Queries.ListTags(c.Request().Context())
	if err != nil {
		slog.Error("list tags", "error", err)
		return c.String(500, "internal server error")
	}
	return c.JSON(200, data)
}

// POST /api/v1/tags
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

// DELETE /api/v1/tags/:title
func (s *Server) deleteTag(c echo.Context) error {
	title := c.Param("title")
	if err := s.Queries.DeleteTag(context.Background(), title); err != nil {
		slog.Error("delete tag by title", "error", err)
		return c.String(500, "internal server error")
	}
	return c.NoContent(204)
}
