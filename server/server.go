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
	api.GET("/photos", s.getAllPhotos)                                                          // list all photos (GET /api/v1/photos)
	api.GET("/photos/:id", s.getPhotoByIDHandler())                                             // get photo by ID (GET /api/v1/photos/:id)
	api.POST("/photos", s.addPhotoHandler(), RequireAdminMiddleware)                            // add photo (POST /api/v1/photos) - admin only
	api.PATCH("/photos/:id/tag/:title", s.addTagToPhotoHandler(), RequireAdminMiddleware)       // add tag to photo (POST /api/v1/photos/tag) - admin only
	api.DELETE("/photos/:id/tag/:title", s.removeTagFromPhotoHandler(), RequireAdminMiddleware) // remove tag from photo (DELETE /api/v1/photos/tag/:title) - admin only

	// bucket (photo storage) endpoints (/api/v1/bucket)
	api.GET("/objects", s.listAllBucketPhotoObjects, RequireAdminMiddleware)             // list all bucket objects (GET /api/v1/bucket) - admin only
	api.POST("/objects", s.uploadPhotoToBucketHandler(), RequireAdminMiddleware)         // upload photo to bucket (POST /api/v1/bucket) - admin only
	api.PATCH("/objects/:name", s.updateObjectMetadataHandler(), RequireAdminMiddleware) // update object metadata (PATCH /api/v1/bucket/object) - admin only

	api.GET("/posts", s.listPosts, IsAdminMiddleware)                                           // list all posts (GET /api/v1/posts) -- admins see all, others see only published
	api.GET("/posts/:slug", s.getPostBySlug, IsAdminMiddleware)                                 // get post by slug (GET /api/v1/posts/:slug) -- admins can see unpublished posts
	api.POST("/posts", s.addPostHandler(), RequireAdminMiddleware)                              // add post (POST /api/v1/posts) - admin only
	api.PATCH("/posts/:slug/tag/:title", s.addTagToPostHandler(), RequireAdminMiddleware)       // add tag to post (POST /api/v1/posts/tag) - admin only
	api.DELETE("/posts/:slug/tag/:title", s.removeTagFromPostHandler(), RequireAdminMiddleware) // remove tag from post (DELETE /api/v1/posts/tag/:title) - admin only

	// tags endpoints (/api/v1/tags)
	api.GET("/tags", s.listTags)                                    // list all tags (GET /api/v1/tags)
	api.POST("/tags", s.addTagHandler(), RequireAdminMiddleware)    // add tag (POST /api/v1/tags) - admin only
	api.DELETE("/tags/:title", s.deleteTag, RequireAdminMiddleware) // delete tag (DELETE /api/v1/tags) - admin only

	// admin endpoints (/api/v1/admin)
	api.GET("/admin", s.isAdmin)                                                                   // check if admin (GET /api/v1/admin)
	api.POST("/admin", s.adminLoginHandler(), NewRateLimiter(5, 10*time.Minute, false).Middleware) // admin login (POST /api/v1/admin) - uses global 5 requests per 10 minutes rate limiter

	return e.Start(env.DefaultEnv.ADDR)
}
