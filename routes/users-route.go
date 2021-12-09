package routes

import (
	"github.com/beeerlian/go-mongo/controllers" // replace
	"github.com/gofiber/fiber/v2"
)

func UsersRoute(route fiber.Router) {
	route.Post("/register/", controllers.UserRegistration)
	route.Post("/login/email/", controllers.LoginWithEmail)
	route.Delete("/:id", controllers.DeleteUser)
	route.Get("/", controllers.GetAllUser)
	route.Post("/participant/:eventId/:userId", controllers.JoinEvent)
	route.Get("/:id", controllers.GetUser)
}
