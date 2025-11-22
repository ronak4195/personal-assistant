package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ronak4195/personal-assistant/internal/handlers"
)

type Handlers struct {
	AuthHandler        *handlers.AuthHandler
	CategoryHandler    *handlers.CategoryHandler
	TransactionHandler *handlers.TransactionHandler
	ReportHandler      *handlers.ReportHandler
	ReminderHandler    *handlers.ReminderHandler
}

func RegisterV1Routes(e *echo.Echo, h Handlers, jwtMiddleware echo.MiddlewareFunc) {
	healthHandler := handlers.NewHealthHandler()

	v1 := e.Group("/api/v1")

	// Health
	v1.GET("/health", healthHandler.Health)

	// Auth (public)
	auth := v1.Group("/auth")
	auth.POST("/signup", h.AuthHandler.Signup)
	auth.POST("/login", h.AuthHandler.Login)
	auth.GET("/me", h.AuthHandler.Me, jwtMiddleware)

	// Protected
	api := v1.Group("", jwtMiddleware)

	// Categories
	api.POST("/categories", h.CategoryHandler.Create)
	api.GET("/categories", h.CategoryHandler.List)
	api.GET("/categories/:id", h.CategoryHandler.Get)
	api.PUT("/categories/:id", h.CategoryHandler.Update)
	api.DELETE("/categories/:id", h.CategoryHandler.Delete)

	// Transactions
	api.POST("/transactions", h.TransactionHandler.Create)
	api.GET("/transactions", h.TransactionHandler.List)
	api.GET("/transactions/:id", h.TransactionHandler.Get)
	api.PUT("/transactions/:id", h.TransactionHandler.Update)
	api.DELETE("/transactions/:id", h.TransactionHandler.Delete)

	// Reports
	api.GET("/reports/summary", h.ReportHandler.Summary)

	// Reminders
	api.POST("/reminders", h.ReminderHandler.Create)
	api.GET("/reminders", h.ReminderHandler.List)
	api.GET("/reminders/:id", h.ReminderHandler.Get)
	api.PUT("/reminders/:id", h.ReminderHandler.Update)
	api.DELETE("/reminders/:id", h.ReminderHandler.Delete)
}
