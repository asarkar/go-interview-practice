package books

import "errors"

// PartialBook holds the updatable fields of a book (pointers distinguish "not provided" from "empty")
type PartialBook struct {
	Title         *string `json:"title"`
	Author        *string `json:"author"`
	PublishedYear *int    `json:"published_year"`
	ISBN          *string `json:"isbn"`
	Description   *string `json:"description"`
}

// NewPartialBook returns a PartialBook. Empty strings and zero are omitted (for partial updates).
func NewPartialBook(title, author string, publishedYear int, isbn, description string) PartialBook {
	p := PartialBook{}
	if title != "" {
		p.Title = strPtr(title)
	}
	if author != "" {
		p.Author = strPtr(author)
	}
	if publishedYear != 0 {
		p.PublishedYear = intPtr(publishedYear)
	}
	if isbn != "" {
		p.ISBN = strPtr(isbn)
	}
	if description != "" {
		p.Description = strPtr(description)
	}
	return p
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

// Book represents a book in the database
type Book struct {
	PartialBook

	ID string `json:"id" gorm:"primaryKey"`
}

// BookRepository defines the operations for book data access
type BookRepository interface {
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Create(book *Book) error
	Update(id string, book *Book) error
	Patch(id string, updates *PartialBook) (*Book, error)
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
	PartiallyUpdateBook(id string, updates *PartialBook) (*Book, error)
	DeleteBook(id string) error
	SearchBooksByAuthor(author string) ([]*Book, error)
	SearchBooksByTitle(title string) ([]*Book, error)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

var errNotFound = errors.New("not found")

type validationError struct{ msg string }

func (e *validationError) Error() string { return e.msg }
