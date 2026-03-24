package user

import "server/pkg/response"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	response.Response
	Token   string `json:"token,omitempty"`
	IsAdmin bool   `json:"isAdmin"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Captcha  string `json:"captcha"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	response.Response
	Token   string `json:"token,omitempty"`
	IsAdmin bool   `json:"isAdmin"`
}

type CaptchaRequest struct {
	Email string `json:"email" binding:"required"`
}

type CaptchaResponse struct {
	response.Response
}
