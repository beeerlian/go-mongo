package main

import (
	"log"
	"os"

	"github.com/beeerlian/go-mongo/config"
	"github.com/beeerlian/go-mongo/controllers"
	"github.com/beeerlian/go-mongo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":     true,
			"message":     "You are at the root endpoint 😉",
			"github_repo": "https://github.com/beeerlian/go-mongo",
		})
	})
	app.Get("/register/", controllers.UserRegistration)

	api := app.Group("/api")

	routes.EventsRoute(api.Group("/events"))
	routes.UsersRoute(api.Group("/users"))
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	config.ConnectDB()

	setupRoutes(app)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}
