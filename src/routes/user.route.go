package routes

import (
	"github.com/DivyanshTiwari001/cjbackend/src/controllers"
	"github.com/DivyanshTiwari001/cjbackend/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/health-check",func(c *fiber.Ctx)error{
		return c.JSON(fiber.Map{"message":"server is healthy"})
	})

	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	protected := api.Group("/protected", middlewares.AuthMiddleware)
	protected.Get("/get-user", controllers.GetUser)
	protected.Get("/logout",controllers.Logout)
}
