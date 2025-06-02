package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/computer101/go-fiber-api/handlers"
)

func main() {
	app := fiber.New()

	// Basic health-check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	// Item routes
	api := app.Group("/items")
	api.Get("/", handlers.ListItems)     // GET  /items
	api.Post("/", handlers.CreateItem)   // POST /items
	api.Get("/:id", handlers.GetItem)    // GET  /items/:id
	api.Put("/:id", handlers.UpdateItem) // PUT  /items/:id
	api.Delete("/:id", handlers.DeleteItem) // DELETE /items/:id

	// Start server on port 8080
	app.Listen(":8080")
}
