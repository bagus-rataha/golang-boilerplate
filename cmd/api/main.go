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
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "fiber-api-boilerplate/docs" // Swagger docs
)

// @title Fiber API Boilerplate
// @version 1.0
// @description Simple REST API with Fiber, GORM, JWT
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.LoadConfig()

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
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.Logger())

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

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
