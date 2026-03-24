package main

import (
	"fmt"
	"log"

	"server/infra/cache"
	"server/infra/config"
	"server/infra/db"
	"server/infra/mq"
	"server/internal/chat"
	"server/internal/router"
	"server/internal/session"
	"server/internal/user"
)

func startServer(addr string, port int) error {
	r := router.New()
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

func main() {
	cfg := config.GetConfig()

	if err := db.InitMysql(); err != nil {
		log.Println("InitMysql error,", err)
		return
	}

	if err := db.Migrate(new(user.User), new(session.Session), new(chat.Message)); err != nil {
		log.Println("migrate error,", err)
		return
	}

	userService := user.NewService(user.NewRepository())
	if err := userService.EnsureConfiguredAdmin(); err != nil {
		log.Println("ensure admin error,", err)
		return
	}

	chatRepo := chat.NewRepository()
	if err := chatRepo.EnsureMessageIdempotency(); err != nil {
		log.Println("ensure message idempotency error,", err)
		return
	}

	cache.Init()
	log.Println("redis init success")

	if err := mq.InitRabbitMQ(); err != nil {
		log.Println("rabbitmq init degraded mode:", err)
	} else {
		log.Println("rabbitmq init success")
	}

	if err := startServer(cfg.Host, cfg.Port); err != nil {
		panic(err)
	}
}
