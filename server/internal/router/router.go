package router

import (
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
	r.Use(gin.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestMetrics())

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
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
	group.POST("/captcha", handler.HandleCaptcha)
}

func registerAIRoutes(group *gin.RouterGroup, sessionHandler *session.Handler, chatHandler *chat.Handler) {
	group.GET("/chat/sessions", sessionHandler.GetUserSessionsByUserName)
	group.POST("/chat/session/rename", sessionHandler.RenameSession)
	group.POST("/chat/session/pin", sessionHandler.UpdateSessionPin)
	group.POST("/chat/session/archive", sessionHandler.UpdateSessionArchive)
	group.POST("/chat/send-new-session", chatHandler.CreateSessionAndSendMessage)
	group.POST("/chat/send", chatHandler.ChatSend)
	group.POST("/chat/history", chatHandler.ChatHistory)
	group.POST("/chat/send-stream-new-session", chatHandler.CreateStreamSessionAndSendMessage)
	group.POST("/chat/send-stream", chatHandler.ChatStreamSend)
}

func registerAdminRoutes(group *gin.RouterGroup, handler *admin.Handler) {
	group.GET("/metrics/all", handler.AllMetrics)
}
