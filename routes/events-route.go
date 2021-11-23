package routes

import (
	"github.com/beeerlian/go-mongo/controllers" // replace
	"github.com/gofiber/fiber/v2"
)

func EventsRoute(route fiber.Router) {
	route.Get("/", controllers.GetAllEvents)
	route.Get("/:id", controllers.GetEvent)
	route.Post("/", controllers.AddEvent)
	route.Put("/:id", controllers.UpdateEvent)
	route.Delete("/:id", controllers.DeleteEvent)
}
