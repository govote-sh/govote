package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/govote-sh/govote/internal/secrets"
	"github.com/govote-sh/govote/internal/tui"

	_ "golang.org/x/crypto/x509roots/fallback"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	err := secrets.SetupSecrets()
	if err != nil {
		log.Error("Could not set up secrets", "error", err)
	}

	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath("/data/govote"),
		wish.WithMiddleware(
			bubbletea.Middleware(tui.TeaHandler),
			logging.Middleware(),
		),
		wish.WithIdleTimeout(8*time.Minute),
		wish.WithMaxTimeout(60*time.Minute),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	// BUG: Timeout is not working
	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
