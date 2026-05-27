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

func Init() error {
	cfg := config.GetConfig()
	addr := cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if _, err := Rdb.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("redis connect failed: %w", err)
	}
	return nil
}

func SetCaptchaForEmail(email, captcha string) error {
	key := generateCaptchaKey(email)
	return Rdb.Set(ctx, key, captcha, 2*time.Minute).Err()
}

// CheckCaptchaForEmail 检查邮箱验证码
func CheckCaptchaForEmail(email, userInput string) (bool, error) {
	// 1. 从缓存中获取验证码
	key := generateCaptchaKey(email)
	// 2. 检查验证码是否存在
	storedCaptcha, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		// 3. 检查错误是否为Nil，即验证码不存在
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
