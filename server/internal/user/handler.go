package user

import (
	"context"

	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Login(ctx context.Context, username, rawPassword string) (string, bool, error)
	Register(ctx context.Context, email, rawPassword, captcha string) (string, bool, error)
	SendCaptcha(ctx context.Context, email string) error
}

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	req := new(LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	token, isAdmin, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}

	res := &LoginResponse{
		Token:   token,
		IsAdmin: isAdmin,
	}
	response.OK(c, res)
}

func (h *Handler) Register(c *gin.Context) {
	req := new(RegisterRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	token, isAdmin, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Captcha)
	if err != nil {
		response.Fail(c, err)
		return
	}

	res := &RegisterResponse{
		Token:   token,
		IsAdmin: isAdmin,
	}
	response.OK(c, res)
}

func (h *Handler) HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, apperror.Wrap(code.CodeInvalidParams, err, code.CodeInvalidParams.Msg()))
		return
	}

	if err := h.service.SendCaptcha(c.Request.Context(), req.Email); err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, &CaptchaResponse{})
}
