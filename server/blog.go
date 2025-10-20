package server

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
)

// GET /api/v1/posts
func (s *Server) listPosts(c echo.Context) error {
	if c.Get("is_admin").(bool) {
		posts, err := s.Queries.ListPostsWithTags(c.Request().Context())
		if err != nil {
			slog.Error("list posts with tags", "error", err)
			return err
		}
		return c.JSON(http.StatusOK, posts)
	} else {
		posts, err := s.Queries.ListPublishedPostsWithTags(c.Request().Context())
		if err != nil {
			slog.Error("list published posts with tags", "error", err)
			return err
		}
		return c.JSON(http.StatusOK, posts)
	}
}

// GET /api/v1/posts/:slug
func (s *Server) getPostBySlug(c echo.Context) error {
	slug := c.Param("slug")
	post, err := s.Queries.GetPostBySlugWithTags(c.Request().Context(), slug)
	if err != nil {
		slog.Error("get post by slug with tags", "error", err)
		return c.JSON(404, echo.Map{"error": "post not found"})
	}
	if !post.Published && !c.Get("is_admin").(bool) { // not published and not admin -- forbidden
		return c.JSON(404, echo.Map{"error": "post not found"})
	}
	return c.JSON(http.StatusOK, post) // either published or user is admin
}

// POST /api/v1/posts
func (s *Server) addPostHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Slug      string   `json:"slug"`
		Published bool     `json:"published"`
		Content   string   `json:"content"`
		Tags      []string `json:"tags" required:"false"`
	}) error {
		tx, err := s.Conn.Begin(c.Request().Context())
		if err != nil {
			slog.Error("begin transaction", "error", err)
			return c.String(500, "internal server error")
		}
		defer tx.Rollback(c.Request().Context())
		queries := s.Queries.WithTx(tx)
		err = queries.CreatePost(c.Request().Context(), db.CreatePostParams{
			Slug:      req.Slug,
			Published: req.Published,
			Content:   req.Content,
		})
		if err != nil {
			return c.JSON(400, echo.Map{"error": "unable to create post"})
		}
		if len(req.Tags) > 0 {
			err = queries.AddTagsToPost(c.Request().Context(), db.AddTagsToPostParams{
				PostSlug: req.Slug,
				Column2:  req.Tags,
			})
			if err != nil {
				slog.Error("add tags to post", "error", err)
				return c.String(500, "internal server error")
			}
		}
		if err := tx.Commit(c.Request().Context()); err != nil {
			slog.Error("commit transaction", "error", err)
			return c.String(500, "internal server error")
		}
		return c.NoContent(http.StatusCreated)
	})
}

func (s *Server) addTagToPostHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Slug  string `param:"slug"`
		Title string `param:"title"`
	}) error {
		err := s.Queries.AddTagToPost(c.Request().Context(), db.AddTagToPostParams{
			PostSlug: req.Slug,
			TagTitle: req.Title,
		})
		if err != nil {
			slog.Error("add tag to post", "error", err)
			return c.String(404, "post or tag not found")
		}
		return c.NoContent(204)
	})
}

func (s *Server) removeTagFromPostHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Slug  string `param:"slug"`
		Title string `param:"title"`
	}) error {
		err := s.Queries.RemoveTagFromPost(c.Request().Context(), db.RemoveTagFromPostParams{
			PostSlug: req.Slug,
			TagTitle: req.Title,
		})
		if err != nil {
			slog.Error("remove tag from post", "error", err)
			return c.String(404, "post or tag not found")
		}
		return c.NoContent(204)
	})
}
