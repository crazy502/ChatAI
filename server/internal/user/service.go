package user

import (
	"context"
	"strings"
	"time"

	"server/infra/cache"
	"server/infra/config"
	"server/infra/mail"
	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/jwt"
	"server/pkg/password"
	"server/pkg/utils"

	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin@qq.com"
	defaultAdminEmail    = "admin@qq.com"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(ctx context.Context, identifier, rawPassword string) (string, bool, error) {
	email := normalizeEmail(identifier)
	if err := checkIdentityRateLimit(ctx, "login-email", email, 10, 15*time.Minute); err != nil {
		return "", false, err
	}

	userInfo, err := s.repo.GetByEmail(email)
	if err == gorm.ErrRecordNotFound {
		return "", false, apperror.New(code.CodeUserNotExist, code.CodeUserNotExist.Msg()).
			WithField("email", email)
	}
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "query user by email failed").
			WithField("email", email)
	}

	if !password.CheckPassword(userInfo.Password, rawPassword) {
		return "", false, apperror.New(code.CodeInvalidPassword, code.CodeInvalidPassword.Msg()).
			WithField("email", email)
	}
	_ = cache.ResetRateLimit(ctx, "login-email", email)

	token, err := jwt.GenerateToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "generate login token failed").
			WithField("user_id", userInfo.ID)
	}

	return token, userInfo.IsAdmin, nil
}

func (s *Service) Register(ctx context.Context, email, rawPassword, captcha string) (string, bool, error) {
	email = normalizeEmail(email)
	if err := checkIdentityRateLimit(ctx, "register-email", email, 5, time.Hour); err != nil {
		return "", false, err
	}

	_, err := s.repo.GetByEmailConsistent(email)
	if err == nil {
		return "", false, apperror.New(code.CodeEmailExist, code.CodeEmailExist.Msg()).
			WithField("email", email)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "query email failed").
			WithField("email", email)
	}

	ok, err := cache.CheckCaptchaForEmail(email, captcha)
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "verify captcha failed").
			WithField("email", email)
	}
	if !ok {
		return "", false, apperror.New(code.CodeInvalidCaptcha, code.CodeInvalidCaptcha.Msg()).
			WithField("email", email)
	}

	username := email
	hashedPassword, err := password.HashPassword(rawPassword)
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "generate password hash failed").
			WithField("email", email)
	}

	userInfo, err := s.repo.Create(username, email, hashedPassword, false)
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "create user failed").
			WithField("email", email).
			WithField("username", username)
	}

	token, err := jwt.GenerateToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "generate register token failed").
			WithField("user_id", userInfo.ID)
	}

	return token, userInfo.IsAdmin, nil
}

func (s *Service) SendCaptcha(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	if err := checkIdentityRateLimit(ctx, "captcha-email", email, 3, 10*time.Minute); err != nil {
		return err
	}

	// 1. 生成随机验证码
	sendCode, err := utils.GetRandomNumbers(6)
	if err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "generate captcha failed").
			WithField("email", email)
	}
	// 2. 存储验证码到缓存
	if err := cache.SetCaptchaForEmail(email, sendCode); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "store captcha failed").
			WithField("email", email)
	}

	// 3. 发送验证码到邮箱
	if err := mail.SendCaptcha(email, sendCode, mail.CodeMsg); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "send captcha email failed").
			WithField("email", email)
	}

	return nil
}

func (s *Service) EnsureConfiguredAdmin() error {
	cfg := config.GetConfig()

	adminUsername := strings.TrimSpace(cfg.AdminConfig.Username)
	if adminUsername == "" {
		adminUsername = defaultAdminUsername
	}

	adminPassword := strings.TrimSpace(cfg.AdminConfig.Password)
	adminEmail := strings.TrimSpace(cfg.AdminConfig.Email)
	if adminEmail == "" {
		adminEmail = defaultAdminEmail
	}
	adminEmail = normalizeEmail(adminEmail)

	passwordHash := ""
	if adminPassword != "" {
		hashed, err := password.HashPassword(adminPassword)
		if err != nil {
			return apperror.Wrap(code.CodeServerBusy, err, "initialize admin password hash failed")
		}
		passwordHash = hashed
	}

	if err := s.repo.EnsureConfiguredAdmin(adminUsername, adminEmail, passwordHash); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "initialize admin account failed").
			WithField("username", adminUsername).
			WithField("email", adminEmail)
	}

	return nil
}

// normalizeEmail 规范化邮箱地址
func normalizeEmail(email string) string {
	// 1. 转换为小写并移除首尾空格
	email = strings.ToLower(strings.TrimSpace(email))
	// 2. 检查邮箱是否为空
	if email == "" {
		return ""
	}
	// 3. 检查邮箱是否包含@符号
	if !strings.Contains(email, "@") {
		return email + "@qq.com"
	}
	return email
}

func checkIdentityRateLimit(ctx context.Context, namespace, identity string, limit int64, window time.Duration) error {
	allowed, _, err := cache.AllowRequest(ctx, namespace, identity, limit, window)
	if err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "check identity rate limit failed")
	}
	if !allowed {
		return apperror.New(code.CodeTooManyRequests, code.CodeTooManyRequests.Msg())
	}
	return nil
}
