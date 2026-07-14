package cache

import (
	"context"
	"crypto/sha256"
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

var rateLimitScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
return {count, ttl}
`)

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

func Close() error {
	if Rdb == nil {
		return nil
	}
	err := Rdb.Close()
	Rdb = nil
	return err
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

func AllowRequest(requestCtx context.Context, namespace, identity string, limit int64, window time.Duration) (bool, time.Duration, error) {
	if Rdb == nil {
		return false, 0, fmt.Errorf("redis client is not initialized")
	}
	if requestCtx == nil {
		requestCtx = ctx
	}
	if namespace == "" || identity == "" || limit <= 0 || window <= 0 {
		return false, 0, fmt.Errorf("invalid rate limit arguments")
	}

	key := rateLimitKey(namespace, identity)
	result, err := rateLimitScript.Run(requestCtx, Rdb, []string{key}, window.Milliseconds()).Result()
	if err != nil {
		return false, 0, err
	}
	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return false, 0, fmt.Errorf("unexpected rate limit result")
	}
	count, err := redisResultInt64(values[0])
	if err != nil {
		return false, 0, err
	}
	ttlMillis, err := redisResultInt64(values[1])
	if err != nil {
		return false, 0, err
	}
	if ttlMillis < 0 {
		ttlMillis = window.Milliseconds()
	}
	return count <= limit, time.Duration(ttlMillis) * time.Millisecond, nil
}

func ResetRateLimit(requestCtx context.Context, namespace, identity string) error {
	if Rdb == nil || namespace == "" || identity == "" {
		return nil
	}
	if requestCtx == nil {
		requestCtx = ctx
	}
	return Rdb.Del(requestCtx, rateLimitKey(namespace, identity)).Err()
}

func rateLimitKey(namespace, identity string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(identity))))
	return fmt.Sprintf("rate:%s:%x", namespace, sum[:])
}

func redisResultInt64(value interface{}) (int64, error) {
	switch typed := value.(type) {
	case int64:
		return typed, nil
	case string:
		return strconv.ParseInt(typed, 10, 64)
	case []byte:
		return strconv.ParseInt(string(typed), 10, 64)
	default:
		return 0, fmt.Errorf("unexpected redis integer type %T", value)
	}
}
