package routes

import (
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/jchen/go-sql-practice/internal/frontend"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
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
		e.App.Logger().Error("Frontend assets not available", "error", err)
		return
	}

	e.App.Logger().Info("Initializing static file serving")

	spaFilesystem := &spaFS{fs: frontendFS}

	e.Router.Bind(&hook.Handler[*core.RequestEvent]{
		Func: func(ev *core.RequestEvent) error {
			// Skip API routes and admin routes
			path := ev.Request.URL.Path
			if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/_/") {
				return ev.Next()
			}

			// Try to serve the static file
			file, err := spaFilesystem.Open(strings.TrimPrefix(path, "/"))
			if err != nil {
				// File not found, continue to next handler
				return ev.Next()
			}
			defer file.Close()

			// Get file info for content type detection
			stat, err := file.Stat()
			if err != nil {
				return ev.Next()
			}

			// Convert fs.File to io.ReadSeeker if possible
			if readSeeker, ok := file.(io.ReadSeeker); ok {
				http.ServeContent(ev.Response, ev.Request, stat.Name(), stat.ModTime(), readSeeker)
				return nil
			}

			// If not a ReadSeeker, continue to next handler
			return ev.Next()
		},
		Priority: -500, // Run before other route handlers but after auth
	})

	e.App.Logger().Info("Static file serving enabled")
}