package main

import (
	"log"
	"MTG/server/database"
	"MTG/server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.ConnectDB()

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/api/items", func(c *fiber.Ctx) error {

		var items []models.Item
		database.DB.Find(&items)
		return c.Status(200).JSON(items)
	})

	log.Fatal(app.Listen(":8000"))
}
