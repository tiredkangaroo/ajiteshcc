package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
)

func (s *Server) getAllPhotos(c echo.Context) error {
	data, err := s.Queries.GetAllPhotosWithTags(c.Request().Context())
	if err != nil {
		slog.Error("get all photos with tags", "error", err)
		return c.String(500, "internal server error")
	}
	return c.JSON(200, data)
}

func (s *Server) getPhotoByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		slog.Error("parse photo id", "error", err)
		return c.String(400, "bad request: invalid photo id")
	}
	data, err := s.Queries.GetPhotoByIDWithTags(c.Request().Context(), int32(id))
	if err != nil {
		slog.Error("get photo by id", "error", err)
		return c.String(500, "internal server error")
	}
	return c.JSON(200, data)
}

func (s *Server) addPhotoHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
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
		if err := s.Queries.AddPhoto(context.Background(), db.AddPhotoParams{
			Title:    pgText(req.Title),
			PhotoUrl: req.PhotoURL,
			Comment:  pgText(req.Comment),
			Metadata: metadata,
		}); err != nil {
			slog.Error("add photo", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	})
}
