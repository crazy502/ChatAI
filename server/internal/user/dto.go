package user

import "server/pkg/response"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=254,endswith=@qq.com"`
	Password string `json:"password" binding:"required,max=72"`
}

type LoginResponse struct {
	response.Response
	Token   string `json:"token,omitempty"`
	IsAdmin bool   `json:"isAdmin"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=254,endswith=@qq.com"`
	Captcha  string `json:"captcha" binding:"required,len=6,numeric"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type RegisterResponse struct {
	response.Response
	Token   string `json:"token,omitempty"`
	IsAdmin bool   `json:"isAdmin"`
}

type CaptchaRequest struct {
	Email string `json:"email" binding:"required,email,max=254,endswith=@qq.com"`
}

type CaptchaResponse struct {
	response.Response
}
