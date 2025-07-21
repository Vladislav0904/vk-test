package items

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ItemHandler struct {
	db *gorm.DB
}

func NewItemHandler(db *gorm.DB) *ItemHandler {
	return &ItemHandler{db: db}
}

func RegisterNewRoutes(app *fiber.App, db *gorm.DB, authMiddleware fiber.Handler, optionalAuth fiber.Handler) {
	h := NewItemHandler(db)
	app.Post("/item", authMiddleware, h.CreateItem)
	app.Get("/items", optionalAuth, h.GetItem)
}
