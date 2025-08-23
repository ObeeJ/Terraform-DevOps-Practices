package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "CarbonAPI v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Initialize database and cache
	db := initDatabase()
	cache := initRedis()
	defer db.Close()
	defer cache.Close()

	// Initialize services
	carbonService := NewCarbonService(db, cache)

	// Routes
	setupRoutes(app, carbonService)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "CarbonAPI",
			"version": "1.0.0",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🌱 CarbonAPI starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App, carbonService *CarbonService) {
	api := app.Group("/api/v1")

	// Carbon calculation endpoints
	api.Post("/calculate", carbonService.CalculateCarbon)
	api.Get("/activities", carbonService.GetActivities)
	api.Get("/factors", carbonService.GetEmissionFactors)
	api.Get("/analytics", carbonService.GetAnalytics)

	// Documentation
	api.Get("/docs", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "CarbonAPI Documentation",
			"endpoints": map[string]interface{}{
				"POST /api/v1/calculate": "Calculate carbon footprint for an activity",
				"GET /api/v1/activities": "List all supported activities",
				"GET /api/v1/factors":    "Get emission factors database",
				"GET /api/v1/analytics":  "Usage analytics and statistics",
				"GET /health":            "Health check endpoint",
			},
			"example": map[string]interface{}{
				"url":    "POST /api/v1/calculate",
				"method": "POST",
				"body": map[string]interface{}{
					"activity":  "shipping",
					"weight":    500,
					"from":      "NYC",
					"to":        "London",
					"transport": "air",
				},
			},
		})
	})
}
