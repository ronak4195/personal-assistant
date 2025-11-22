package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ronak4195/personal-assistant/internal/config"
	"github.com/ronak4195/personal-assistant/internal/db"
	"github.com/ronak4195/personal-assistant/internal/handlers"
	appmw "github.com/ronak4195/personal-assistant/internal/middleware"
	"github.com/ronak4195/personal-assistant/internal/repositories"
	"github.com/ronak4195/personal-assistant/internal/routes"
	"github.com/ronak4195/personal-assistant/internal/services"

	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load env / config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, database, err := db.Connect(ctx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer func() {
		_ = dbClient.Disconnect(context.Background())
	}()

	// rm
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Global middleware
	e.Use(echomw.Recover())
	e.Use(echomw.Logger())
	e.Use(appmw.RequestID())

	// 🔑 CORS
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{
			cfg.FrontEndURL,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	// Repositories
	userRepo := repositories.NewUserRepository(database)
	categoryRepo := repositories.NewCategoryRepository(database)
	transactionRepo := repositories.NewTransactionRepository(database)
	reminderRepo := repositories.NewReminderRepository(database)

	// Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo, categoryRepo)
	reportService := services.NewReportService(transactionRepo, categoryRepo)
	reminderService := services.NewReminderService(reminderRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)
	reminderHandler := handlers.NewReminderHandler(reminderService)

	// JWT middleware
	jwtMiddleware := appmw.JWTAuth(cfg.JWTSecret)

	// Routes
	routes.RegisterV1Routes(e, routes.Handlers{
		AuthHandler:        authHandler,
		CategoryHandler:    categoryHandler,
		TransactionHandler: transactionHandler,
		ReportHandler:      reportHandler,
		ReminderHandler:    reminderHandler,
	}, jwtMiddleware)

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(":" + cfg.HTTPPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := e.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exiting")
}
