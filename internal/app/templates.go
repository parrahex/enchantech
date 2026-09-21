package app

import (
	"html/template"
	"log/slog"
	"path/filepath"
)

func loadTemplates(contentDir string) *template.Template {
	pattern := filepath.Join(contentDir, "templates", "*.html")

	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		slog.Warn("no page templates found, serving placeholder", "pattern", pattern)

		return nil
	}

	parsed, err := template.ParseGlob(pattern)
	if err != nil {
		slog.Error("page templates could not be parsed, serving placeholder", "error", err)

		return nil
	}

	return parsed
}
