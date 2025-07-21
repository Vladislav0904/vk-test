package items

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strings"
	"time"
)

type createItemRequest struct {
	Title    string  `json:"title"`
	Text     string  `json:"text"`
	ImageURL string  `json:"image_url"`
	Price    float32 `json:"price"`
}

type itemResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url"`
	Price     float32   `json:"price"`
	AuthorID  uint64    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

func validateItem(req createItemRequest) (string, bool) {
	if len(req.Title) == 0 || len(req.Title) > 100 {
		return "invalid title", false
	}
	if len(req.Text) == 0 || len(req.Text) > 1000 {
		return "invalid text", false
	}
	if len(req.ImageURL) > 255 {
		return "invalid image_url", false
	}
	if req.ImageURL != "" {
		if _, err := url.ParseRequestURI(req.ImageURL); err != nil {
			return "invalid image_url", false
		}
	}
	if req.Price <= 0 || req.Price > 1_000_000 {
		return "invalid price", false
	}
	return "", true
}

func (h *ItemHandler) CreateItem(ctx *fiber.Ctx) error {
	var req createItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	req.Title = strings.TrimSpace(req.Title)
	req.Text = strings.TrimSpace(req.Text)
	req.ImageURL = strings.TrimSpace(req.ImageURL)
	if msg, ok := validateItem(req); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}
	userID, ok := ctx.Locals("user_id").(uint64)
	if !ok {
		f, fok := ctx.Locals("user_id").(float64)
		if !fok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		userID = uint64(f)
	}
	item := Item{
		Title:     req.Title,
		Text:      req.Text,
		ImageURL:  req.ImageURL,
		Price:     req.Price,
		AuthorID:  userID,
		CreatedAt: time.Now(),
	}
	if err := h.db.Create(&item).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(itemResponse{
		ID:        item.ID,
		Title:     item.Title,
		Text:      item.Text,
		ImageURL:  item.ImageURL,
		Price:     item.Price,
		AuthorID:  item.AuthorID,
		CreatedAt: item.CreatedAt,
	})
}

func (h *ItemHandler) GetItem(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	pageSize := ctx.QueryInt("page_size", 10)
	sortBy := ctx.Query("sort_by", "created_at")
	order := ctx.Query("order", "desc")
	minPrice := ctx.QueryFloat("min_price", 0)
	maxPrice := ctx.QueryFloat("max_price", 0)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	dbq := h.db.Model(&Item{})
	if minPrice > 0 {
		dbq = dbq.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		dbq = dbq.Where("price <= ?", maxPrice)
	}
	if sortBy != "created_at" && sortBy != "price" {
		sortBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	dbq = dbq.Order(sortBy + " " + order)
	offset := (page - 1) * pageSize
	var items []Item
	if err := dbq.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	var currentUserID uint64
	if v := ctx.Locals("user_id"); v != nil {
		if id, ok := v.(uint64); ok {
			currentUserID = id
		} else if f, ok := v.(float64); ok {
			currentUserID = uint64(f)
		}
	}
	type feedItem struct {
		ID        uint64    `json:"id"`
		Title     string    `json:"title"`
		Text      string    `json:"text"`
		ImageURL  string    `json:"image_url"`
		Price     float32   `json:"price"`
		AuthorID  *uint64   `json:"author_id,omitempty"`
		CreatedAt time.Time `json:"created_at"`
		IsMine    bool      `json:"is_mine,omitempty"`
	}
	resp := make([]feedItem, 0, len(items))
	for _, it := range items {
		fi := feedItem{
			ID:        it.ID,
			Title:     it.Title,
			Text:      it.Text,
			ImageURL:  it.ImageURL,
			Price:     it.Price,
			CreatedAt: it.CreatedAt,
		}
		if currentUserID != 0 {
			fi.AuthorID = &it.AuthorID
			if it.AuthorID == currentUserID {
				fi.IsMine = true
			}
		}
		resp = append(resp, fi)
	}
	return ctx.JSON(resp)
}
