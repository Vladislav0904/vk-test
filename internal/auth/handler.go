package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func RegisterNewRouter(app *fiber.App, db *gorm.DB) {
	h := NewAuthHandler(db)
	app.Post("/login", h.Login)
	app.Post("/register", h.Register)
}
