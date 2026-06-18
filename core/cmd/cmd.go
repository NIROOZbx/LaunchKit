package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Launchkit-org/LaunchKit/core/internal/app"
)

func Run(a *app.App) error {
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP}
	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		a.Logger.Info().Str("addr", a.Listener.Addr().String()).Msg("core grpc server listening")
		if err := a.Server.Serve(a.Listener); err != nil {
			errChan <- fmt.Errorf("grpc serve: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		a.Logger.Info().Msg("shutdown signal received")
	}

	done := make(chan struct{})
	go func() {
		a.Server.GracefulStop()  // blocks until all in-flight RPCs finish
		close(done)
	}()

	select {
	case <-done:
		a.Logger.Info().Msg("grpc server stopped gracefully")
	case <-time.After(10 * time.Second):
		a.Logger.Warn().Msg("graceful stop timed out, forcing stop")
		a.Server.Stop()
	}

	a.DBPool.Close()
	a.Logger.Info().Msg("shutdown complete")
	return nil
}