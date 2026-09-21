package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/parrahex/enchantech/internal/config"
)

// App owns the HTTP lifecycle of the application.
type App struct {
	web *http.Server
}

// New creates the Gin router and configures the HTTP runtime.
func New(settings config.Config) *App {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	templates := loadTemplates(settings.ContentDir)
	if templates != nil {
		router.SetHTMLTemplate(templates)
	}

	if directory, ok := assetDirectory(settings.ContentDir); ok {
		router.Static("/assets", directory)
	}

	registerRoutes(router, templates != nil)

	return &App{
		web: &http.Server{
			Addr:              settings.Address,
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// Run blocks until the HTTP runtime stops or returns an unexpected error.
func (application *App) Run() error {
	if err := application.web.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Shutdown stops accepting new requests and waits for active requests until
// the supplied context expires.
func (application *App) Shutdown(ctx context.Context) error {
	return application.web.Shutdown(ctx)
}
