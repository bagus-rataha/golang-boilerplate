package middleware

import (
	"fiber-api-boilerplate/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupMiddleware configures all global middleware
func SetupMiddleware(app *fiber.App, cfg *config.Config) {
	// Panic recovery
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
	}))

	// CORS
	setupCORS(app, cfg)

	// Rate limiting (production only)
	if cfg.IsProduction() {
		setupRateLimiter(app, cfg)
	}

	// Compression (production only)
	if cfg.IsProduction() {
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))
	}

	// Request logger
	app.Use(Logger())
}

// setupCORS configures CORS middleware
func setupCORS(app *fiber.App, cfg *config.Config) {
	if cfg.IsProduction() {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.AllowedOrigins,
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: true,
		}))
	} else {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		}))
	}
}

// setupRateLimiter configures rate limiting
func setupRateLimiter(app *fiber.App, cfg *config.Config) {
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: cfg.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too many requests, please try again later",
				"data":    nil,
			})
		},
	}))
}
