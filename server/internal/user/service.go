package user

import (
	"log"
	"strings"

	"server/infra/cache"
	"server/infra/config"
	"server/infra/mail"
	"server/pkg/code"
	"server/pkg/jwt"
	"server/pkg/password"
	"server/pkg/utils"

	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "admin"
	defaultAdminEmail    = "admin@gopherai.local"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(username, rawPassword string) (string, bool, code.Code) {
	userInfo, err := s.repo.GetByUsername(username)
	if err == gorm.ErrRecordNotFound {
		return "", false, code.CodeUserNotExist
	}
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	if password.IsBcryptHash(userInfo.Password) {
		if !password.CheckPassword(userInfo.Password, rawPassword) {
			return "", false, code.CodeInvalidPassword
		}
	} else {
		if userInfo.Password != utils.MD5(rawPassword) {
			return "", false, code.CodeInvalidPassword
		}

		hashedPassword, err := password.HashPassword(rawPassword)
		if err != nil {
			log.Println("login hash legacy password error:", err)
		} else if err := s.repo.UpdatePassword(userInfo.ID, hashedPassword); err != nil {
			log.Println("login migrate legacy password error:", err)
		}
	}

	token, err := jwt.GenerateToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	return token, userInfo.IsAdmin, code.CodeSuccess
}

func (s *Service) Register(email, rawPassword, captcha string) (string, bool, code.Code) {
	_, err := s.repo.GetByEmail(email)
	if err == nil {
		return "", false, code.CodeEmailExist
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", false, code.CodeServerBusy
	}

	ok, err := cache.CheckCaptchaForEmail(email, captcha)
	if err != nil {
		return "", false, code.CodeServerBusy
	}
	if !ok {
		return "", false, code.CodeInvalidCaptcha
	}

	username := utils.GetRandomNumbers(11)
	hashedPassword, err := password.HashPassword(rawPassword)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	userInfo, err := s.repo.Create(username, email, hashedPassword, false)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	if err := mail.SendCaptcha(email, username, mail.UserNameMsg); err != nil {
		return "", false, code.CodeServerBusy
	}

	token, err := jwt.GenerateToken(userInfo.ID, userInfo.Username, userInfo.IsAdmin)
	if err != nil {
		return "", false, code.CodeServerBusy
	}

	return token, userInfo.IsAdmin, code.CodeSuccess
}

func (s *Service) SendCaptcha(email string) code.Code {
	sendCode := utils.GetRandomNumbers(6)
	if err := cache.SetCaptchaForEmail(email, sendCode); err != nil {
		return code.CodeServerBusy
	}

	if err := mail.SendCaptcha(email, sendCode, mail.CodeMsg); err != nil {
		return code.CodeServerBusy
	}

	return code.CodeSuccess
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

	passwordHash, err := password.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	return s.repo.EnsureConfiguredAdmin(adminUsername, adminEmail, passwordHash)
}
