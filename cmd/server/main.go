package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"exit-debt/internal/config"
	"exit-debt/internal/database"
	"exit-debt/internal/handler"
	"exit-debt/internal/middleware"
	"exit-debt/internal/service"
)

func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Initialize database with GORM
	db, err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize GORM services
	authService, err := service.NewAuthServiceGORM(db.DB, cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize auth service")
	}

	contactService := service.NewContactServiceGORM(db.DB)
	debtService := service.NewDebtServiceGORM(db.DB)

	// Initialize GORM handlers
	authHandler := handler.NewAuthHandlerGORM(authService)
	contactHandler := handler.NewContactHandlerGORM(contactService)
	debtHandler := handler.NewDebtHandlerGORM(debtService)

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Public routes
	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddlewareGORM(authService))
	{
		// Health check
		protected.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Contact routes
		protected.POST("/contacts", contactHandler.CreateContact)
		protected.GET("/contacts", contactHandler.GetUserContacts)
		protected.GET("/contacts/:id", contactHandler.GetContact)
		protected.PUT("/contacts/:id", contactHandler.UpdateContact)
		protected.DELETE("/contacts/:id", contactHandler.DeleteContact)

		// Debt List routes
		protected.POST("/debt-lists", debtHandler.CreateDebtList)
		protected.GET("/debt-lists", debtHandler.GetUserDebtLists)
		protected.GET("/debt-lists/:id", debtHandler.GetDebtList)
		protected.PUT("/debt-lists/:id", debtHandler.UpdateDebtList)
		protected.DELETE("/debt-lists/:id", debtHandler.DeleteDebtList)

		// Debt Item routes
		protected.POST("/debt-items", debtHandler.CreateDebtItem)
		protected.GET("/debt-items/:id", debtHandler.GetDebtItem)
		protected.PUT("/debt-items/:id", debtHandler.UpdateDebtItem)
		protected.DELETE("/debt-items/:id", debtHandler.DeleteDebtItem)
		protected.GET("/debt-lists/:id/items", debtHandler.GetDebtListItems)

		// Special routes
		protected.GET("/debt-items/overdue", debtHandler.GetOverdueItems)
		protected.GET("/debt-items/due-soon", debtHandler.GetDueSoonItems)
		protected.GET("/debt-lists/:id/payment-schedule", debtHandler.GetPaymentSchedule)
		protected.GET("/debt-lists/:id/payments", debtHandler.GetTotalPaymentsForDebtList)
		protected.GET("/upcoming-payments", debtHandler.GetUpcomingPayments)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Info().Msgf("Starting server on %s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
} 