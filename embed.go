package main

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

func getFrontendFS() (fs.FS, error) {
	subFS, err := fs.Sub(frontendFiles, "frontend/dist")
	if err != nil {
		return nil, fmt.Errorf("failed to create frontend sub-filesystem: %w", err)
	}
	return subFS, nil
}

func frontendExists() bool {
	_, err := getFrontendFS()
	return err == nil
}