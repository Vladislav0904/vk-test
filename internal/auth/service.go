package auth

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"vk-test/pkg/utils"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func validateUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

func validatePassword(password string) bool {
	return len(password) >= 6 && len(password) <= 64
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req registerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	req.Username = strings.TrimSpace(req.Username)
	if !validateUsername(req.Username) || !validatePassword(req.Password) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid username or password format"})
	}
	var exists int64
	if err := h.db.Model(&User{}).Where("username = ?", req.Username).Count(&exists).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	if exists > 0 {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "username already exists"})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "hash error"})
	}
	user := User{Username: req.Username, Password: string(hash)}
	if err := h.db.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(userResponse{ID: user.ID, Username: user.Username})
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req loginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	req.Username = strings.TrimSpace(req.Username)
	if !validateUsername(req.Username) || !validatePassword(req.Password) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid username or password format"})
	}
	var user User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return ctx.JSON(tokenResponse{Token: token})
}
