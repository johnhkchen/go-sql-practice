package frontend

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed ../../frontend/dist/*
var frontendFiles embed.FS

func GetFrontendFS() (fs.FS, error) {
	subFS, err := fs.Sub(frontendFiles, "frontend/dist/client")
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend sub-filesystem: %w", err)
	}
	return subFS, nil
}

func FrontendExists() bool {
	_, err := GetFrontendFS()
	return err == nil
}