package main

import (
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/handlers"
	"fiber-api-boilerplate/internal/middleware"
	"fiber-api-boilerplate/internal/repository"
	"fiber-api-boilerplate/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "fiber-api-boilerplate/docs" // Swagger docs
)

// @title Fiber API Boilerplate
// @version 2.0
// @description Production-ready REST API with Fiber, GORM, JWT, Validation
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Log environment
	log.Printf("Environment: %s", cfg.AppEnv)

	// Connect to database
	db := config.ConnectDB(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		Prefork:      cfg.IsProduction(),
		ServerHeader: "",
	})

	// Global middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
	}))

	// CORS configuration
	if cfg.IsProduction() && len(cfg.AllowedOrigins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     joinOrigins(cfg.AllowedOrigins),
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: true,
		}))
	} else {
		app.Use(cors.New())
	}

	// Rate limiting (production only)
	if cfg.IsProduction() {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimitMax,
			Expiration: cfg.RateLimitWindow,
		}))
	}

	app.Use(middleware.Logger())

	// Swagger documentation (development only)
	if cfg.IsDevelopment() {
		app.Get("/swagger/*", swagger.HandlerDefault)
		log.Printf("Swagger UI: http://localhost:%s/swagger/", cfg.Port)
	}

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	users := api.Group("/users")
	users.Use(middleware.JWTProtected(cfg.JWTSecret))
	users.Get("/me", userHandler.GetProfile)
	users.Put("/me", userHandler.UpdateProfile)
	users.Get("/", userHandler.ListUsers) // Admin only in real case

	// Start server
	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

// joinOrigins joins allowed origins into a comma-separated string
func joinOrigins(origins []string) string {
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ","
		}
		result += origin
	}
	return result
}
