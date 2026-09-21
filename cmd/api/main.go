package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/parrahex/enchantech/internal/app"
	"github.com/parrahex/enchantech/internal/config"
)

func main() {
	// Text logs go to stdout so the hosting environment or a log collector can
	// collect them without coupling the application to a logging platform.
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)

	if err := startApplication(); err != nil {
		slog.Error("application exited with an error", "error", err)
		os.Exit(1)
	}
}

func startApplication() error {
	settings, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	application := app.New(settings)

	// Run blocks while the application accepts requests, so it runs in a
	// goroutine while the main flow waits for a signal.
	applicationErrors := make(chan error, 1)

	go func() {
		applicationErrors <- application.Run()
	}()

	return gracefulShutdown(application, applicationErrors, settings.ShutdownTimeout)
}

func gracefulShutdown(application *app.App, applicationErrors <-chan error, timeout time.Duration) error {
	signalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-applicationErrors:
		// A result before a shutdown signal usually means startup failed, for
		// example because the address could not be opened.
		if err == nil {
			return fmt.Errorf("application stopped unexpectedly")
		}

		return fmt.Errorf("application run failed: %w", err)

	case <-signalContext.Done():
		// Managed shutdown normally sends SIGTERM. Ctrl+C sends SIGINT.
		slog.Info("shutdown signal received")

		// A separate context gives active requests their own shutdown deadline.
		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			timeout,
		)
		defer cancel()

		if err := application.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("shutdown application: %w", err)
		}

		// Wait until Run has returned before allowing the process to exit.
		if err := <-applicationErrors; err != nil {
			return fmt.Errorf("finish application shutdown: %w", err)
		}

		return nil
	}
}
