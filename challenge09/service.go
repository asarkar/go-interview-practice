package books

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
	if book.Title == nil || *book.Title == "" || book.Author == nil || *book.Author == "" {
		return &validationError{"title and author are required"}
	}
	return s.repo.Create(book)
}

func (s *DefaultBookService) UpdateBook(id string, book *Book) error {
	if book.Title == nil || *book.Title == "" || book.Author == nil || *book.Author == "" {
		return &validationError{"title and author are required"}
	}
	return s.repo.Update(id, book)
}

func (s *DefaultBookService) PartiallyUpdateBook(id string, updates *PartialBook) (*Book, error) {
	if updates.Title != nil && *updates.Title == "" {
		return nil, &validationError{"title cannot be empty"}
	}
	if updates.Author != nil && *updates.Author == "" {
		return nil, &validationError{"author cannot be empty"}
	}
	return s.repo.Patch(id, updates)
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
