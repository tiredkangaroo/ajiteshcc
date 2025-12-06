package server

import (
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
		if s.musicLastTrack == nil {
			return c.JSON(http.StatusNoContent, map[string]string{"error": "no track currently playing"})
		}
		return c.JSON(http.StatusOK, s.musicLastTrack)
	}
	s.musicLastUpdated = time.Now() // update last checked time

	if s.spotifyHTTPClient == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "spotify unavailable"})
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	resp, err := s.spotifyHTTPClient.Do(req)
	if err != nil {
		slog.Error("spotify now playing request", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return c.JSON(http.StatusNoContent, map[string]string{"error": "no track currently playing"})
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("now playing status code", "status_code", resp.StatusCode)
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "spotify service error"})
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	if result.Item.Type != "track" {
		return c.JSON(http.StatusNoContent, map[string]string{"error": "no track currently playing"})
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
	s.musicLastTrack = &track
	return c.JSON(http.StatusOK, track)
}

func (s *Server) spotifyLoginHandler(c echo.Context) error {
	url := s.spotifyOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
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
	s.spotifyHTTPClient = s.spotifyOAuthConfig.Client(c.Request().Context(), token)
	return c.JSON(http.StatusOK, map[string]string{"status": "spotify login successful"})
}
