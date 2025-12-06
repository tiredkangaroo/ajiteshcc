package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

const musicStaleDuration = 15 * time.Second

type Track struct {
	Name     string `json:"name"`
	Artists  string `json:"artists"`
	CoverURL string `json:"cover_url"`
}

func (s *Server) getNowPlaying(c echo.Context) error {
	// avoid spamming spotify api if we recently checked
	if time.Since(s.musicLastUpdated) < musicStaleDuration {
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}
	s.musicLastUpdated = time.Now() // update last checked time

	if s.spotifyHTTPClient == nil {
		s.musicLastResponseCode = http.StatusServiceUnavailable
		s.musicLastResponseBody = map[string]string{"error": "spotify unavailable"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		slog.Error("create spotify now playing request", "error", err)
		s.musicLastResponseCode = http.StatusInternalServerError
		s.musicLastResponseBody = map[string]string{"error": "internal server error"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}
	resp, err := s.spotifyHTTPClient.Do(req)
	if err != nil {
		slog.Error("spotify now playing request", "error", err)
		s.musicLastResponseCode = http.StatusInternalServerError
		s.musicLastResponseBody = map[string]string{"error": "internal server error"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		s.musicLastResponseCode = http.StatusNoContent
		s.musicLastResponseBody = map[string]string{"error": "no track currently playing"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("now playing status code", "status_code", resp.StatusCode)
		s.musicLastResponseCode = http.StatusBadGateway
		s.musicLastResponseBody = map[string]string{"error": "spotify service error"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}

	var result struct {
		Item struct {
			Name    string `json:"name"`
			Type    string `json:"type"` // "track", "episode", "ad", just keeping it just in case i either go insane enough to listen to a podcast on spotify or spotify starts giving me ads on premium...
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		slog.Error("decode spotify now playing response", "error", err)
		s.musicLastResponseCode = http.StatusInternalServerError
		s.musicLastResponseBody = map[string]string{"error": "internal server error"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}

	if result.Item.Type != "track" {
		s.musicLastResponseCode = http.StatusNoContent
		s.musicLastResponseBody = map[string]string{"error": "no track currently playing"}
		return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
	}

	artists := make([]string, len(result.Item.Artists))
	for i, artist := range result.Item.Artists {
		artists[i] = artist.Name
	}

	imageURL := ""
	if len(result.Item.Album.Images) > 0 {
		imageURL = result.Item.Album.Images[0].URL
	}

	track := Track{
		Name:     result.Item.Name,
		Artists:  strings.Join(artists, ", "),
		CoverURL: imageURL,
	}
	s.musicLastResponseCode = http.StatusOK
	s.musicLastResponseBody = track
	return c.JSON(s.musicLastResponseCode, s.musicLastResponseBody)
}

func (s *Server) spotifyLoginHandler(c echo.Context) error {
	url := s.spotifyOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

func (s *Server) spotifyLogoutHandler(c echo.Context) error {
	s.spotifyToken = nil
	s.spotifyHTTPClient = nil
	s.musicLastUpdated = time.Time{}
	s.musicLastResponseCode = 0
	s.musicLastResponseBody = nil
	return redirectFrontend(c, "/")
}

func (s *Server) spotifyCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing code in callback"})
	}
	token, err := s.spotifyOAuthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		slog.Error("spotify oauth exchange", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to exchange code for token"})
	}
	s.spotifyToken = token
	s.spotifyHTTPClient = s.spotifyOAuthConfig.Client(context.Background(), token)
	s.musicLastUpdated = time.Time{}
	s.musicLastResponseCode = 0
	s.musicLastResponseBody = nil
	return redirectFrontend(c, "/")
}
