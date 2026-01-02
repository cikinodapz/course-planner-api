package handler

import (
	"course-planner-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	service service.BookService
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(s service.BookService) *BookHandler {
	return &BookHandler{service: s}
}

// SearchBooks godoc
// @Summary Search books from Google Books API
// @Description Search for books by query string. Consumes Google Books API with API Key authentication.
// @Tags Books (External API)
// @Security BearerAuth
// @Param query query string true "Search query (e.g., 'algoritma', 'pemrograman')"
// @Param maxResults query int false "Maximum results to return (default: 10, max: 40)"
// @Success 200 {object} service.BookSearchResponse
// @Failure 400 {object} map[string]string "Bad Request - query is required"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/books [get]
func (h *BookHandler) SearchBooks(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter is required",
		})
	}

	maxResults := 10
	if maxResultsStr := c.Query("maxResults"); maxResultsStr != "" {
		if parsed, err := strconv.Atoi(maxResultsStr); err == nil {
			maxResults = parsed
		}
	}

	result, err := h.service.SearchBooks(query, maxResults)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

// GetBookDetail godoc
// @Summary Get book detail from Google Books API
// @Description Get detailed information about a specific book by its Google Books volume ID
// @Tags Books (External API)
// @Security BearerAuth
// @Param id path string true "Google Books Volume ID"
// @Success 200 {object} service.BookDetail
// @Failure 400 {object} map[string]string "Bad Request - invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/books/{id} [get]
func (h *BookHandler) GetBookDetail(c *fiber.Ctx) error {
	volumeID := c.Params("id")
	if volumeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "book ID is required",
		})
	}

	book, err := h.service.GetBookByID(volumeID)
	if err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(book)
}
