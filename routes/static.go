package routes

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/jchen/go-sql-practice/internal/frontend"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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
	frontendFS, err := frontend.GetFrontendFS()
	if err != nil {
		e.App.Logger().Error("Failed to get frontend filesystem", "error", err)
		return
	}

	spaFilesystem := &spaFS{fs: frontendFS}
	e.Router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "/",
		Filesystem: http.FS(spaFilesystem),
		Browse:     false,
	}))
}