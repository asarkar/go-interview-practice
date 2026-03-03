package books

import "errors"

// DefaultBookService implements BookService
type DefaultBookService struct {
	repo BookRepository
}

// NewBookService creates a new book service
func NewBookService(repo BookRepository) *DefaultBookService {
	return &DefaultBookService{repo: repo}
}

func (s *DefaultBookService) GetAllBooks() ([]*Book, error) {
	return s.repo.GetAll()
}

func (s *DefaultBookService) GetBookByID(id string) (*Book, error) {
	return s.repo.GetByID(id)
}

func (s *DefaultBookService) CreateBook(book *Book) error {
	if book.Title == "" || book.Author == "" {
		return errors.New("title and author are required")
	}
	return s.repo.Create(book)
}

func (s *DefaultBookService) UpdateBook(id string, book *Book) error {
	if book.Title == "" || book.Author == "" {
		return errors.New("title and author are required")
	}
	return s.repo.Update(id, book)
}

func (s *DefaultBookService) PartiallyUpdateBook(id string, patch *BookPatch) (*Book, error) {
	if patch.Title != nil && *patch.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if patch.Author != nil && *patch.Author == "" {
		return nil, errors.New("author cannot be empty")
	}
	return s.repo.Patch(id, patch)
}

func (s *DefaultBookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

func (s *DefaultBookService) SearchBooksByAuthor(author string) ([]*Book, error) {
	return s.repo.SearchByAuthor(author)
}

func (s *DefaultBookService) SearchBooksByTitle(title string) ([]*Book, error) {
	return s.repo.SearchByTitle(title)
}
