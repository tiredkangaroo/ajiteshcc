package server

import (
	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/bucket"
	"github.com/tiredkangaroo/ajiteshcc/env"
)

func (s *Server) listAllBucketPhotoObjects(c echo.Context) error {
	objects, err := bucket.ListAllObjectsInBucket("photos")
	if err != nil {
		return c.JSON(500, map[string]string{"error": "failed to list objects"})
	}
	var objectsWithPublicURL []map[string]string
	pub_url := *env.DefaultEnv.R2_PHOTOS_BUCKET_PUBLIC_URL
	for _, obj := range objects {
		pub_url.Path = obj.Key // set the path to the object key
		objectsWithPublicURL = append(objectsWithPublicURL, map[string]string{
			"name": obj.Key,
			"url":  pub_url.String(),
		})
	}
	return c.JSON(200, objectsWithPublicURL)
}
