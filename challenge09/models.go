package books

import "errors"

// Book represents a book in the database
type Book struct {
	ID            string `json:"id"             gorm:"primaryKey"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
	ISBN          string `json:"isbn"`
	Description   string `json:"description"`
}

// BookPatch represents optional fields for partial updates
type BookPatch struct {
	Title         *string `json:"title"`
	Author        *string `json:"author"`
	PublishedYear *int    `json:"published_year"`
	ISBN          *string `json:"isbn"`
	Description   *string `json:"description"`
}

// BookRepository defines the operations for book data access
type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Patch(id string, patch *BookPatch) (*Book, error)
	Delete(id string) error
	SearchByAuthor(author string) ([]*Book, error)
	SearchByTitle(title string) ([]*Book, error)
}

// BookService defines the business logic for book operations
type BookService interface {
	GetAllBooks() ([]*Book, error)
	GetBookByID(id string) (*Book, error)
	CreateBook(book *Book) error
	UpdateBook(id string, book *Book) error
	PartiallyUpdateBook(id string, patch *BookPatch) (*Book, error)
	DeleteBook(id string) error
	SearchBooksByAuthor(author string) ([]*Book, error)
	SearchBooksByTitle(title string) ([]*Book, error)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

var errNotFound = errors.New("not found")
