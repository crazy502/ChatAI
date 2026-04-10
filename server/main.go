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
		observe.Error(ctx, "初始化 MySQL 失败", err)
		os.Exit(1)
	}

	if err := db.Migrate(new(user.User), new(session.Session), new(chat.Message)); err != nil {
		observe.Error(ctx, "执行数据库迁移失败", err)
		os.Exit(1)
	}

	userService := user.NewService(user.NewRepository())
	if err := userService.EnsureConfiguredAdmin(); err != nil {
		observe.Error(ctx, "初始化管理员账号失败", err)
		os.Exit(1)
	}

	chatRepo := chat.NewRepository()
	if err := chatRepo.EnsureMessageIdempotency(); err != nil {
		observe.Error(ctx, "初始化消息幂等索引失败", err)
		os.Exit(1)
	}

	if err := cache.Init(); err != nil {
		observe.Error(ctx, "初始化 Redis 失败", err)
		os.Exit(1)
	}
	observe.Info(ctx, "redis init success")

	if err := mq.InitRabbitMQ(); err != nil {
		observe.Warn(ctx, "rabbitmq init degraded mode", "cause", err.Error())
	} else {
		observe.Info(ctx, "rabbitmq init success")
	}

	if err := startServer(cfg.Host, cfg.Port); err != nil {
		observe.Error(ctx, "启动 HTTP 服务失败", err)
		os.Exit(1)
	}
}
