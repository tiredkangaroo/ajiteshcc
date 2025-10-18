package server

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/env"
	"github.com/tiredkangaroo/ajiteshcc/gen/db"
)

type Server struct {
	Conn    *pgx.Conn
	Queries *db.Queries
}

func (s *Server) Run() error {
	e := echo.New()

	api := e.Group("/api/v1")

	// photos endpoints (/api/v1/photos)
	api.GET("/photos", s.getAllPhotos)                               // list all photos (GET /api/v1/photos)
	api.GET("/photos/:id", s.getPhotoByID)                           // get photo by ID (GET /api/v1/photos/:id)
	api.POST("/photos", s.addPhotoHandler(), RequireAdminMiddleware) // add photo (POST /api/v1/photos) - admin only

	// bucket (photo storage) endpoints (/api/v1/bucket)
	api.GET("/bucket", s.listAllBucketPhotoObjects) // list all bucket objects (GET /api/v1/bucket)

	// tags endpoints (/api/v1/tags)
	api.GET("/tags", s.listTags)                                      // list all tags (GET /api/v1/tags)
	api.POST("/tags", s.addTagHandler(), RequireAdminMiddleware)      // add tag (POST /api/v1/tags) - admin only
	api.DELETE("/tags", s.deleteTagHandler(), RequireAdminMiddleware) // delete tag (DELETE /api/v1/tags) - admin only

	// admin endpoints (/api/v1/admin)
	api.GET("/admin", s.isAdmin)                                                                   // check if admin (GET /api/v1/admin)
	api.POST("/admin", s.adminLoginHandler(), NewRateLimiter(5, 10*time.Minute, false).Middleware) // admin login (POST /api/v1/admin) - uses global 5 requests per 10 minutes rate limiter

	return e.Start(env.DefaultEnv.ADDR)
}
