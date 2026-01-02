package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// BookService interface for Google Books API operations
type BookService interface {
	SearchBooks(query string, maxResults int) (*BookSearchResponse, error)
	GetBookByID(volumeID string) (*BookDetail, error)
}

// bookService implements BookService
type bookService struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewBookService creates a new BookService instance
func NewBookService() BookService {
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	return &bookService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://www.googleapis.com/books/v1/volumes",
	}
}

// BookSearchResponse represents the search response
type BookSearchResponse struct {
	TotalItems int         `json:"total_items"`
	Books      []BookBrief `json:"books"`
}

// BookBrief represents a brief book item in search results
type BookBrief struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
	Description   string   `json:"description,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	InfoLink      string   `json:"info_link,omitempty"`
	ISBN10        string   `json:"isbn_10,omitempty"`
	ISBN13        string   `json:"isbn_13,omitempty"`
}

// BookDetail represents detailed book information
type BookDetail struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle,omitempty"`
	Authors       []string `json:"authors,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
	Description   string   `json:"description,omitempty"`
	PageCount     int      `json:"page_count,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	AverageRating float64  `json:"average_rating,omitempty"`
	RatingsCount  int      `json:"ratings_count,omitempty"`
	Language      string   `json:"language,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	PreviewLink   string   `json:"preview_link,omitempty"`
	InfoLink      string   `json:"info_link,omitempty"`
	ISBN10        string   `json:"isbn_10,omitempty"`
	ISBN13        string   `json:"isbn_13,omitempty"`
}

// Google Books API response structures
type googleBooksSearchResponse struct {
	TotalItems int               `json:"totalItems"`
	Items      []googleBooksItem `json:"items"`
}

type googleBooksItem struct {
	ID         string           `json:"id"`
	VolumeInfo googleVolumeInfo `json:"volumeInfo"`
}

type googleVolumeInfo struct {
	Title               string                     `json:"title"`
	Subtitle            string                     `json:"subtitle"`
	Authors             []string                   `json:"authors"`
	Publisher           string                     `json:"publisher"`
	PublishedDate       string                     `json:"publishedDate"`
	Description         string                     `json:"description"`
	PageCount           int                        `json:"pageCount"`
	Categories          []string                   `json:"categories"`
	AverageRating       float64                    `json:"averageRating"`
	RatingsCount        int                        `json:"ratingsCount"`
	Language            string                     `json:"language"`
	ImageLinks          *googleImageLinks          `json:"imageLinks"`
	PreviewLink         string                     `json:"previewLink"`
	InfoLink            string                     `json:"infoLink"`
	IndustryIdentifiers []googleIndustryIdentifier `json:"industryIdentifiers"`
}

type googleImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

type googleIndustryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// SearchBooks searches for books using Google Books API
func (s *bookService) SearchBooks(query string, maxResults int) (*BookSearchResponse, error) {
	if s.apiKey == "" {
		return nil, errors.New("GOOGLE_BOOKS_API_KEY is not configured")
	}

	if query == "" {
		return nil, errors.New("query parameter is required")
	}

	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 40 {
		maxResults = 40 // Google Books API max
	}

	// Build URL with query parameters
	reqURL := fmt.Sprintf("%s?q=%s&maxResults=%d&key=%s",
		s.baseURL,
		url.QueryEscape(query),
		maxResults,
		s.apiKey,
	)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Books API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google Books API returned status: %d", resp.StatusCode)
	}

	var googleResp googleBooksSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Transform to our response format
	books := make([]BookBrief, 0, len(googleResp.Items))
	for _, item := range googleResp.Items {
		book := BookBrief{
			ID:            item.ID,
			Title:         item.VolumeInfo.Title,
			Authors:       item.VolumeInfo.Authors,
			Publisher:     item.VolumeInfo.Publisher,
			PublishedDate: item.VolumeInfo.PublishedDate,
			Description:   truncateDescription(item.VolumeInfo.Description, 200),
			InfoLink:      item.VolumeInfo.InfoLink,
		}

		if item.VolumeInfo.ImageLinks != nil {
			book.Thumbnail = item.VolumeInfo.ImageLinks.Thumbnail
		}

		// Extract ISBNs
		for _, identifier := range item.VolumeInfo.IndustryIdentifiers {
			switch identifier.Type {
			case "ISBN_10":
				book.ISBN10 = identifier.Identifier
			case "ISBN_13":
				book.ISBN13 = identifier.Identifier
			}
		}

		books = append(books, book)
	}

	return &BookSearchResponse{
		TotalItems: googleResp.TotalItems,
		Books:      books,
	}, nil
}

// GetBookByID gets detailed book information by volume ID
func (s *bookService) GetBookByID(volumeID string) (*BookDetail, error) {
	if s.apiKey == "" {
		return nil, errors.New("GOOGLE_BOOKS_API_KEY is not configured")
	}

	if volumeID == "" {
		return nil, errors.New("volume ID is required")
	}

	reqURL := fmt.Sprintf("%s/%s?key=%s", s.baseURL, volumeID, s.apiKey)

	resp, err := s.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Books API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("book not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google Books API returned status: %d", resp.StatusCode)
	}

	var item googleBooksItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	book := &BookDetail{
		ID:            item.ID,
		Title:         item.VolumeInfo.Title,
		Subtitle:      item.VolumeInfo.Subtitle,
		Authors:       item.VolumeInfo.Authors,
		Publisher:     item.VolumeInfo.Publisher,
		PublishedDate: item.VolumeInfo.PublishedDate,
		Description:   item.VolumeInfo.Description,
		PageCount:     item.VolumeInfo.PageCount,
		Categories:    item.VolumeInfo.Categories,
		AverageRating: item.VolumeInfo.AverageRating,
		RatingsCount:  item.VolumeInfo.RatingsCount,
		Language:      item.VolumeInfo.Language,
		PreviewLink:   item.VolumeInfo.PreviewLink,
		InfoLink:      item.VolumeInfo.InfoLink,
	}

	if item.VolumeInfo.ImageLinks != nil {
		book.Thumbnail = item.VolumeInfo.ImageLinks.Thumbnail
	}

	// Extract ISBNs
	for _, identifier := range item.VolumeInfo.IndustryIdentifiers {
		switch identifier.Type {
		case "ISBN_10":
			book.ISBN10 = identifier.Identifier
		case "ISBN_13":
			book.ISBN13 = identifier.Identifier
		}
	}

	return book, nil
}

// truncateDescription truncates description to specified length
func truncateDescription(desc string, maxLen int) string {
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + "..."
}
