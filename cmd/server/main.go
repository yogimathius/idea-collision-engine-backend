package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"idea-collision-engine-api/internal/auth"
	"idea-collision-engine-api/internal/collision"
	"idea-collision-engine-api/internal/database"
	"idea-collision-engine-api/internal/handlers"
	"idea-collision-engine-api/internal/middleware"
	"idea-collision-engine-api/internal/models"
	"idea-collision-engine-api/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connections
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	redis, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	aiService := collision.NewAIService(cfg.OpenAIAPIKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, redis, jwtService)
	collisionHandler := handlers.NewCollisionHandler(db, redis, aiService)
	subscriptionHandler := handlers.NewSubscriptionHandler(db, redis, cfg.StripeSecretKey)

	// Initialize collision engine with domains
	if err := seedCollisionDomains(db); err != nil {
		log.Printf("Warning: Failed to seed collision domains: %v", err)
	}

	if err := collisionHandler.Initialize(); err != nil {
		log.Fatalf("Failed to initialize collision handler: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Idea Collision Engine API",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: errorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins[0],
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"service":   "idea-collision-engine-api",
			"version":   "1.0.0",
			"timestamp": time.Now(),
		})
	})

	// API routes
	api := app.Group("/api")

	// Authentication routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/profile", middleware.AuthMiddleware(jwtService), authHandler.GetProfile)
	auth.Put("/profile", middleware.AuthMiddleware(jwtService), authHandler.UpdateProfile)

	// Collision routes
	collisions := api.Group("/collisions")
	
	// Rate limiting for collision generation
	rateLimitConfig := middleware.RateLimitConfig{
		WindowSeconds: 60,     // 1 minute window
		MaxRequests:   10,     // 10 requests per minute
		SkipPremium:   true,   // Skip rate limiting for premium users
	}
	
	collisions.Post("/generate", 
		middleware.AuthMiddleware(jwtService),
		middleware.UsageLimitMiddleware(db, redis),
		middleware.RateLimitMiddleware(redis, rateLimitConfig),
		collisionHandler.GenerateCollision,
	)
	
	collisions.Get("/history", 
		middleware.AuthMiddleware(jwtService),
		collisionHandler.GetCollisionHistory,
	)
	
	collisions.Put("/:id/rate", 
		middleware.AuthMiddleware(jwtService),
		collisionHandler.RateCollision,
	)
	
	collisions.Get("/usage", 
		middleware.AuthMiddleware(jwtService),
		collisionHandler.GetUsageStatus,
	)
	
	collisions.Get("/health", collisionHandler.HealthCheck)

	// Domain routes
	domains := api.Group("/domains")
	domains.Get("/basic", collisionHandler.GetBasicDomains)
	domains.Get("/premium", 
		middleware.AuthMiddleware(jwtService),
		middleware.RequirePremium(),
		collisionHandler.GetPremiumDomains,
	)

	// Subscription routes
	subscriptions := api.Group("/subscriptions")
	subscriptions.Get("/plans", subscriptionHandler.GetPricingPlans)
	subscriptions.Post("/checkout", 
		middleware.AuthMiddleware(jwtService),
		subscriptionHandler.CreateCheckoutSession,
	)
	subscriptions.Get("/status", 
		middleware.AuthMiddleware(jwtService),
		subscriptionHandler.GetSubscriptionStatus,
	)
	subscriptions.Post("/cancel", 
		middleware.AuthMiddleware(jwtService),
		subscriptionHandler.CancelSubscription,
	)
	subscriptions.Post("/webhook", subscriptionHandler.WebhookHandler)

	// Documentation routes
	docsHandler := handlers.NewDocsHandler()
	docs := app.Group("/docs")
	docs.Get("/", docsHandler.SwaggerUI())
	docs.Get("/openapi.yaml", docsHandler.OpenAPISpec)
	docs.Static("/", "./internal/handlers/swagger-ui")

	// Start server
	port := ":" + cfg.Port
	fmt.Printf("🚀 Idea Collision Engine API starting on port %s\n", cfg.Port)
	fmt.Printf("📊 Environment: %s\n", cfg.Environment)
	
	// Graceful shutdown
	go func() {
		if err := app.Listen(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server stopped")
}

// errorHandler handles application errors
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(models.ErrorResponse{
		Error:   "server_error",
		Message: err.Error(),
		Code:    code,
	})
}

// seedCollisionDomains populates the database with collision domains if empty
func seedCollisionDomains(db *database.PostgresDB) error {
	// Check if domains already exist
	domains, err := db.GetCollisionDomains("basic")
	if err != nil {
		return err
	}

	if len(domains) > 0 {
		fmt.Printf("✅ Collision domains already seeded (%d domains)\n", len(domains))
		return nil
	}

	// Seed the domains
	seedDomains := database.GetCollisionDomainSeeds()
	
	for _, domain := range seedDomains {
		if err := db.CreateCollisionDomain(&domain); err != nil {
			return fmt.Errorf("failed to seed domain %s: %w", domain.Name, err)
		}
	}

	fmt.Printf("✅ Seeded %d collision domains\n", len(seedDomains))
	return nil
}