package handlers

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Item represents our resource.
type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// A thread-safe, in-memory store for Item objects.
var (
	itemStore  = make(map[int]Item)
	nextItemID = 1
	storeMu    sync.Mutex
)

// ListItems returns all items as JSON.
func ListItems(c *fiber.Ctx) error {
	storeMu.Lock()
	defer storeMu.Unlock()

	// Convert map to slice
	items := make([]Item, 0, len(itemStore))
	for _, itm := range itemStore {
		items = append(items, itm)
	}

	return c.JSON(items)
}

// GetItem returns a single item by ID.
func GetItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	item, found := itemStore[id]
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "Item Not Found")
	}
	return c.JSON(item)
}

// CreateItem reads JSON from the request body, creates a new Item, and returns it.
func CreateItem(c *fiber.Ctx) error {
	var newItem Item
	if err := c.BodyParser(&newItem); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}
	if newItem.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Name is required")
	}
	if newItem.Quantity < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Quantity must be non-negative")
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	newItem.ID = nextItemID
	nextItemID++
	itemStore[newItem.ID] = newItem

	c.Status(fiber.StatusCreated)
	return c.JSON(newItem)
}

// UpdateItem updates an existing item partially (name and/or quantity).
func UpdateItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var payload Item
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	existing, found := itemStore[id]
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "Item Not Found")
	}

	// Apply partial updates
	if payload.Name != "" {
		existing.Name = payload.Name
	}
	if payload.Quantity >= 0 {
		existing.Quantity = payload.Quantity
	}
	itemStore[id] = existing

	return c.JSON(existing)
}

// DeleteItem removes an item by ID.
func DeleteItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	if _, found := itemStore[id]; !found {
		return fiber.NewError(fiber.StatusNotFound, "Item Not Found")
	}
	delete(itemStore, id)

	return c.SendStatus(fiber.StatusNoContent) // 204
}
