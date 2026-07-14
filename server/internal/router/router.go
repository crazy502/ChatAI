package router

import (
	"time"

	"server/internal/admin"
	"server/internal/chat"
	"server/internal/middleware"
	"server/internal/session"
	"server/internal/user"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	userRepo := user.NewRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	sessionRepo := session.NewRepository()
	sessionService := session.NewService(sessionRepo)
	sessionHandler := session.NewHandler(sessionService)

	chatRepo := chat.NewRepository()
	chatService := chat.NewService(chatRepo, sessionRepo)
	chat.StartMessageConsumer(chatRepo)
	chatHandler := chat.NewHandler(chatService)

	adminService := admin.NewService()
	adminHandler := admin.NewHandler(adminService)

	r := gin.New()
	r.Use(middleware.RequestContext())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestObserver())
	r.Use(middleware.MaxBodyBytes(1 << 20))

	api := r.Group("/api/v1")

	userGroup := api.Group("/user")
	registerUserRoutes(userGroup, userHandler)

	aiGroup := api.Group("/AI")
	aiGroup.Use(middleware.Auth())
	registerAIRoutes(aiGroup, sessionHandler, chatHandler)

	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.Auth())
	adminGroup.Use(middleware.RequireAdmin())
	registerAdminRoutes(adminGroup, adminHandler)

	return r
}

func registerUserRoutes(group *gin.RouterGroup, handler *user.Handler) {
	group.POST("/register", middleware.RateLimit("register-ip", 10, time.Hour, false), handler.Register)
	group.POST("/login", middleware.RateLimit("login-ip", 30, 15*time.Minute, false), handler.Login)
	group.POST("/captcha", middleware.RateLimit("captcha-ip", 10, 10*time.Minute, false), handler.HandleCaptcha)
}

func registerAIRoutes(group *gin.RouterGroup, sessionHandler *session.Handler, chatHandler *chat.Handler) {
	chatLimiter := middleware.RateLimit("chat-user", 30, time.Minute, true)
	group.GET("/chat/sessions", sessionHandler.GetUserSessionsByUserName)
	group.POST("/chat/session/rename", sessionHandler.RenameSession)
	group.POST("/chat/session/pin", sessionHandler.UpdateSessionPin)
	group.POST("/chat/session/archive", sessionHandler.UpdateSessionArchive)
	group.POST("/chat/send-new-session", chatLimiter, chatHandler.CreateSessionAndSendMessage)
	group.POST("/chat/send", chatLimiter, chatHandler.ChatSend)
	group.POST("/chat/history", chatHandler.ChatHistory)
	group.POST("/chat/send-stream-new-session", chatLimiter, chatHandler.CreateStreamSessionAndSendMessage)
	group.POST("/chat/send-stream", chatLimiter, chatHandler.ChatStreamSend)
}

func registerAdminRoutes(group *gin.RouterGroup, handler *admin.Handler) {
	group.GET("/metrics/all", handler.AllMetrics)
}
