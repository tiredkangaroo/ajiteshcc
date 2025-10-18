package server

import (
	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/bucket"
)

func (s *Server) listAllBucketPhotoObjects(c echo.Context) error {
	objects, err := bucket.ListAllObjectsInBucket("photos")
	if err != nil {
		return c.JSON(500, map[string]string{"error": "failed to list objects"})
	}
	return c.JSON(200, objects)
}
