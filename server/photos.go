package server

import (
	"context"
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
		Tags     []string       `json:"tags" required:"false"`
	}) error {
		tx, err := s.Conn.Begin(c.Request().Context())
		if err != nil {
			slog.Error("begin transaction", "error", err)
			return c.String(500, "internal server error")
		}
		defer tx.Rollback(c.Request().Context())
		queries := s.Queries.WithTx(tx)

		photoID, err := queries.AddPhoto(c.Request().Context(), db.AddPhotoParams{
			Title:    pgText(req.Title),
			PhotoUrl: req.PhotoURL,
			Comment:  pgText(req.Comment),
			Metadata: req.Metadata,
		})
		if err != nil {
			slog.Error("add photo", "error", err)
			return c.String(500, "internal server error")
		}
		if err := queries.AddTagsToPhoto(context.Background(), db.AddTagsToPhotoParams{
			PhotoID: photoID,
			Column2: req.Tags,
		}); err != nil {
			slog.Error("add photo tags", "error", err)
			return c.String(500, "internal server error")
		}

		if err := tx.Commit(c.Request().Context()); err != nil {
			slog.Error("commit transaction", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(204)
	})
}
