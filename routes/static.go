package routes

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

type spaFS struct {
	fs fs.FS
}

func (s *spaFS) Open(name string) (fs.File, error) {
	// Try exact match first
	file, err := s.fs.Open(name)
	if err == nil {
		return file, nil
	}

	// Fallback to index.html for SPA routing
	// Don't handle API routes or admin routes
	if !strings.HasPrefix(name, "api/") && !strings.HasPrefix(name, "_/") {
		return s.fs.Open("index.html")
	}

	return nil, err
}

func registerStatic(e *core.ServeEvent) {
	e.Router.GET("/*", func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Static serving not implemented")
	})
}