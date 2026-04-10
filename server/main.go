package main

import (
	"context"
	"fmt"
	"os"

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

func startServer(addr string, port int) error {
	r := router.New()
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

func main() {
	observe.Init()
	ctx := context.Background()
	cfg := config.GetConfig()

	if err := db.InitMysql(); err != nil {
		observe.Error(ctx, "initialize mysql failed", err)
		os.Exit(1)
	}

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
	observe.Info(ctx, "redis init success")
	observe.Info(ctx, "mysql init success", "read_replica_enabled", db.HasReplica())

	if err := mq.InitRabbitMQ(); err != nil {
		observe.Warn(ctx, "rabbitmq init degraded mode", "cause", err.Error())
	} else {
		observe.Info(ctx, "rabbitmq init success")
	}

	if err := startServer(cfg.Host, cfg.Port); err != nil {
		observe.Error(ctx, "start http server failed", err)
		os.Exit(1)
	}
}
