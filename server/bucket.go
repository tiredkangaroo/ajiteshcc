package server

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta"
	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/bucket"
)

func (s *Server) listAllBucketPhotoObjects(c echo.Context) error {
	objects, err := bucket.ListAllObjectsInBucket(c.Request().Context(), "photos")
	if err != nil {
		slog.Error("list all bucket photo objects", "error", err)
		return c.JSON(500, map[string]string{"error": "unable to list objects"})
	}
	return c.JSON(200, objects)
}

func (s *Server) uploadPhotoToBucketHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Name string `form:"name"` // object key
	}) error {
		fileheader, err := c.FormFile("file")
		if err != nil {
			slog.Error("get uploaded file", "error", err)
			return c.JSON(400, map[string]string{"error": "file is required"})
		}
		file, err := fileheader.Open()
		if err != nil {
			slog.Error("open uploaded file", "error", err)
			return c.JSON(500, map[string]string{"error": "failed to open file"})
		}
		defer file.Close()

		md, err := metadata(file)
		if err != nil {
			slog.Error("extract metadata", "error", err)
			return c.JSON(500, map[string]string{"error": "unable to extract metadata"})
		}

		if seeker, ok := file.(io.ReadSeeker); ok {
			_, err = seeker.Seek(0, io.SeekStart)
			if err != nil {
				slog.Error("reset file reader", "error", err)
				return c.JSON(500, map[string]string{"error": "internal server error"})
			}
		} else {
			slog.Error("file is not seekable")
			return c.JSON(500, map[string]string{"error": "internal server error"})
		}

		// warn: use of req.Name directly can lead to overwriting existing files and
		// security issues. this is admin-only, but req.Name should still be handled.
		if err := bucket.PutObjectInBucket(c.Request().Context(), "photos", req.Name, md, file); err != nil {
			slog.Error("put object in bucket", "error", err)
			return c.JSON(500, map[string]string{"error": "unable to upload file"})
		}
		return c.JSON(200, map[string]string{"success": "file uploaded successfully"})
	})
}

func (s *Server) updateObjectMetadataHandler() echo.HandlerFunc {
	return handler(func(c echo.Context, req struct {
		Metadata map[string]string `json:"metadata"` // new metadata
	}) error {
		name := c.Param("name")
		if err := bucket.UpdateObjectMetadata(c.Request().Context(), "photos", name, req.Metadata); err != nil {
			slog.Error("update object metadata", "error", err)
			return c.JSON(500, map[string]string{"error": "unable to update metadata"})
		}
		return c.JSON(200, map[string]string{"success": "metadata updated successfully"})
	})
}

func metadata(file io.ReadSeeker) (map[string]string, error) {
	exif, err := imagemeta.Decode(file)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"CreatedAt":    exif.CreateDate().String(),
		"LensMake":     exif.LensMake,
		"LensModel":    exif.LensModel,
		"CameraMake":   exif.CameraMake.String(),
		"CameraModel":  exif.CameraModel.String(),
		"Aperture":     fmt.Sprintf("%f", exif.FNumber),
		"FocalLength":  fmt.Sprintf("%f", exif.FocalLength),
		"ISO":          fmt.Sprintf("%d", exif.ISOSpeed),
		"ShutterSpeed": meta.ExposureTime(exif.ExposureTime).String(),
		"ImageType":    exif.ImageType.String(),
		"ImageWidth":   fmt.Sprintf("%d", exif.ImageWidth),
		"ImageHeight":  fmt.Sprintf("%d", exif.ImageHeight),
		"Latitude":     fmt.Sprintf("%f", exif.GPS.Latitude()),
		"Longitude":    fmt.Sprintf("%f", exif.GPS.Longitude()),
		"Altitude":     fmt.Sprintf("%f", exif.GPS.Altitude()),
	}, nil
}
