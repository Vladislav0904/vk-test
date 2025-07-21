package main

import (
	"github.com/gofiber/fiber/v2"
	"vk-test/internal/auth"
	"vk-test/internal/config"
	"vk-test/internal/database"
	"vk-test/internal/items"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	app := fiber.New()
	auth.RegisterNewRouter(app, db)
	items.RegisterNewRoutes(app, db, auth.AuthMiddleware(), auth.OptionalAuthMiddleware())
	app.Listen(":8080")
}
