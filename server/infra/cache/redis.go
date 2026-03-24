package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"server/infra/config"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb = (*redis.Client)(nil)
	ctx = context.Background()
)

func Init() {
	cfg := config.GetConfig()
	addr := cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if _, err := Rdb.Ping(ctx).Result(); err != nil {
		panic("redis connect failed: " + err.Error())
	}
}

func SetCaptchaForEmail(email, captcha string) error {
	key := generateCaptchaKey(email)
	return Rdb.Set(ctx, key, captcha, 2*time.Minute).Err()
}

func CheckCaptchaForEmail(email, userInput string) (bool, error) {
	key := generateCaptchaKey(email)

	storedCaptcha, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if !strings.EqualFold(storedCaptcha, userInput) {
		return false, nil
	}

	if err := Rdb.Del(ctx, key).Err(); err != nil {
		return false, err
	}

	return true, nil
}

func generateCaptchaKey(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
}
