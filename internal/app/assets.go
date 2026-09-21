package app

import (
	"log/slog"
	"os"
	"path/filepath"
)

func assetDirectory(contentDir string) (string, bool) {
	directory := filepath.Join(contentDir, "assets")

	info, err := os.Stat(directory)
	if err != nil || !info.IsDir() {
		slog.Warn("no asset directory found, serving no assets", "path", directory)

		return "", false
	}

	return directory, true
}
