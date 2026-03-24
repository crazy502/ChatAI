package user

import (
	"net/http"

	"server/pkg/code"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	req := new(LoginRequest)
	res := new(LoginResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, isAdmin, resultCode := h.service.Login(req.Username, req.Password)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	res.Token = token
	res.IsAdmin = isAdmin
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Register(c *gin.Context) {
	req := new(RegisterRequest)
	res := new(RegisterResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, isAdmin, resultCode := h.service.Register(req.Email, req.Password, req.Captcha)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	res.Token = token
	res.IsAdmin = isAdmin
	c.JSON(http.StatusOK, res)
}

func (h *Handler) HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	resultCode := h.service.SendCaptcha(req.Email)
	if resultCode != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(resultCode))
		return
	}

	res.Success()
	c.JSON(http.StatusOK, res)
}
