package user

import (
	"context"
	"strings"

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
	defaultAdminPassword = "admin"
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

	token, err := jwt.GenerateToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
	if err != nil {
		return "", false, apperror.Wrap(code.CodeServerBusy, err, "generate login token failed").
			WithField("user_id", userInfo.ID)
	}

	return token, userInfo.IsAdmin, nil
}

func (s *Service) Register(ctx context.Context, email, rawPassword, captcha string) (string, bool, error) {
	email = normalizeEmail(email)

	_, err := s.repo.GetByEmail(email)
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

	sendCode := utils.GetRandomNumbers(6)
	if err := cache.SetCaptchaForEmail(email, sendCode); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "store captcha failed").
			WithField("email", email)
	}

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
	if adminPassword == "" {
		adminPassword = defaultAdminPassword
	}

	adminEmail := strings.TrimSpace(cfg.AdminConfig.Email)
	if adminEmail == "" {
		adminEmail = defaultAdminEmail
	}
	adminEmail = normalizeEmail(adminEmail)

	passwordHash, err := password.HashPassword(adminPassword)
	if err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "initialize admin password hash failed")
	}

	if err := s.repo.EnsureConfiguredAdmin(adminUsername, adminEmail, passwordHash); err != nil {
		return apperror.Wrap(code.CodeServerBusy, err, "initialize admin account failed").
			WithField("username", adminUsername).
			WithField("email", adminEmail)
	}

	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
