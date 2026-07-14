package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/infra/cache"
	"server/infra/config"
	"server/infra/db"
	"server/infra/mq"
	"server/internal/chat"
	"server/internal/router"
	"server/internal/session"
	"server/internal/user"
	"server/pkg/observe"
)

func startServer(ctx context.Context, addr string, port int) error {
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", addr, port),
		Handler:           router.New(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

func main() {
	observe.Init()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := config.InitConfig(); err != nil {
		observe.Error(ctx, "initialize config failed", err)
		os.Exit(1)
	}
	cfg := config.GetConfig()

	if err := db.InitMysql(); err != nil {
		observe.Error(ctx, "initialize mysql failed", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			observe.Warn(context.Background(), "close mysql failed", "cause", err.Error())
		}
	}()

	if err := db.Migrate(new(user.User), new(session.Session), new(chat.Message)); err != nil {
		observe.Error(ctx, "run database migration failed", err)
		os.Exit(1)
	}

	userService := user.NewService(user.NewRepository())
	if err := userService.EnsureConfiguredAdmin(); err != nil {
		observe.Error(ctx, "initialize admin account failed", err)
		os.Exit(1)
	}

	sessionRepo := session.NewRepository()
	if err := sessionRepo.EnsureListIndexes(); err != nil {
		observe.Error(ctx, "initialize session list indexes failed", err)
		os.Exit(1)
	}

	chatRepo := chat.NewRepository()
	if err := chatRepo.EnsureMessageIdempotency(); err != nil {
		observe.Error(ctx, "initialize message idempotency failed", err)
		os.Exit(1)
	}
	if err := chatRepo.EnsureHistoryIndexes(); err != nil {
		observe.Error(ctx, "initialize message history indexes failed", err)
		os.Exit(1)
	}

	if err := cache.Init(); err != nil {
		observe.Error(ctx, "initialize redis failed", err)
		os.Exit(1)
	}
	defer func() {
		if err := cache.Close(); err != nil {
			observe.Warn(context.Background(), "close redis failed", "cause", err.Error())
		}
	}()
	observe.Info(ctx, "redis init success")
	observe.Info(ctx, "mysql init success", "read_replica_enabled", db.HasReplica())

	if err := mq.InitRabbitMQ(); err != nil {
		observe.Warn(ctx, "rabbitmq init degraded mode", "cause", err.Error())
	} else {
		observe.Info(ctx, "rabbitmq init success")
		defer mq.DestroyRabbitMQ()
	}

	if err := startServer(ctx, cfg.Host, cfg.Port); err != nil {
		observe.Error(ctx, "start http server failed", err)
		os.Exit(1)
	}
}
