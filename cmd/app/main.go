package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/esfands/retpaladinbot/config"
	"github.com/esfands/retpaladinbot/internal/bot"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/internal/services/helix"
	"github.com/esfands/retpaladinbot/internal/services/scheduler"
	"github.com/esfands/retpaladinbot/internal/services/turso"
	"golang.org/x/exp/slog"
)

var (
	Version   = "dev"
	Timestamp = "unknown"
)

func main() {
	Timestamp = time.Now().Format(time.RFC3339)

	version := os.Getenv("VERSION")
	if version == "" {
		version = Version
	} else {
		Version = version
	}

	// Intialize the logger depending on the version of the app
	var logger *slog.Logger
	if version == "dev" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	slog.Info("Starting the application", "version", version, "timestamp", Timestamp)

	// Load configuration
	cfg, err := config.New(Version, time.Now())
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	gctx, cancel := global.WithCancel(global.New(context.Background(), cfg))

	{
		slog.Info("Setting up Turso database")
		gctx.Crate().Turso, err = turso.Setup(gctx, turso.SetupOptions{
			URL: cfg.Turso.URL,
		})
		if err != nil {
			slog.Error("Error setting up Turso database", "error", err)
			cancel()
			return
		}
		slog.Info("Turso database setup complete")
	}

	{
		slog.Info("Setting up scheduler")
		gctx.Crate().Scheduler, err = scheduler.Setup(gctx)
		if err != nil {
			slog.Error("Error setting up scheduler", "error", err)
			cancel()
			return
		}

		slog.Info("Scheduler setup complete")
	}

	{
		slog.Info("Setting up Helix API")
		gctx.Crate().Helix, err = helix.Setup(gctx, gctx.Crate().Scheduler, helix.SetupOptions{
			ClientID:     cfg.Twitch.Helix.ClientID,
			ClientSecret: cfg.Twitch.Helix.ClientSecret,
		})
		if err != nil {
			slog.Error("Error setting up Helix API", "error", err)
			cancel()
			return
		}

		slog.Info("Helix API setup complete")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	wg := sync.WaitGroup{}

	go func() {
		<-interrupt
		cancel()

		go func() {
			// If interrupt signal is not handled in 1 minute or interrupted once again, force shutdown
			select {
			case <-time.After(time.Minute):
			case <-interrupt:
			}
			fmt.Println("force shutdown")
		}()

		fmt.Println("shutting down")

		wg.Wait()

		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		bot.StartBot(gctx, cfg, Version)
	}()

	<-done

	slog.Info("Application stopped")
	os.Exit(0)
}
