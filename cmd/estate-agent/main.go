package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonathanpopham/estate-agent/internal/app"
	"github.com/jonathanpopham/estate-agent/internal/config"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		slog.Error("estate-agent failed", "error", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] != "serve" {
		return errors.New("usage: estate-agent serve")
	}

	cfg := config.FromEnv()
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           app.NewServer(cfg),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errs := make(chan error, 1)
	go func() {
		slog.Info("starting estate-agent", "addr", cfg.Addr, "cloud", cfg.Cloud, "dry_run", cfg.DryRun)
		errs <- server.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errs:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve: %w", err)
	}
}
