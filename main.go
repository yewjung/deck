package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Define your API endpoint
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/card", func(c *fiber.Ctx) error {
		deckId := c.Query("deckid")
		_, card, err := DrawFrom(deckId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(card)
	})
	app.Get("/newdeck", func(c *fiber.Ctx) error {
		deckId, err := NewDeck()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendString(deckId)
	})

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
