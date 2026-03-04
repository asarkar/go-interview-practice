package books

import (
	"errors"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORMBookRepository implements BookRepository using SQLite via GORM
type GORMBookRepository struct {
	db *gorm.DB
}

// NewInMemoryBookRepository creates a GORM-backed in-memory SQLite repository
func NewInMemoryBookRepository() *GORMBookRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to open in-memory database: " + err.Error())
	}
	if err := db.AutoMigrate(&Book{}); err != nil {
		panic("failed to migrate schema: " + err.Error())
	}
	return &GORMBookRepository{db: db}
}

func (r *GORMBookRepository) GetAll() ([]*Book, error) {
	var books []*Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *GORMBookRepository) GetByID(id string) (*Book, error) {
	var book Book
	if err := r.db.First(&book, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errNotFound
		}
		return nil, err
	}
	return &book, nil
}

func (r *GORMBookRepository) Create(book *Book) error {
	book.ID = uuid.New().String()
	return r.db.Create(book).Error
}

func (r *GORMBookRepository) Update(id string, book *Book) error {
	py, isbn, desc := 0, "", ""
	if book.PublishedYear != nil {
		py = *book.PublishedYear
	}
	if book.ISBN != nil {
		isbn = *book.ISBN
	}
	if book.Description != nil {
		desc = *book.Description
	}
	result := r.db.Model(&Book{}).Where("id = ?", id).Updates(map[string]any{
		"title":          *book.Title,
		"author":         *book.Author,
		"published_year": py,
		"isbn":           isbn,
		"description":    desc,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errNotFound
	}
	return nil
}

func (r *GORMBookRepository) Patch(id string, updates *PartialBook) (*Book, error) {
	u := make(map[string]any)
	if updates.Title != nil {
		u["title"] = *updates.Title
	}
	if updates.Author != nil {
		u["author"] = *updates.Author
	}
	if updates.PublishedYear != nil {
		u["published_year"] = *updates.PublishedYear
	}
	if updates.ISBN != nil {
		u["isbn"] = *updates.ISBN
	}
	if updates.Description != nil {
		u["description"] = *updates.Description
	}
	if len(u) == 0 {
		return r.GetByID(id)
	}
	result := r.db.Model(&Book{}).Where("id = ?", id).Updates(u)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errNotFound
	}
	return r.GetByID(id)
}

func (r *GORMBookRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&Book{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errNotFound
	}
	return nil
}

func (r *GORMBookRepository) SearchByAuthor(author string) ([]*Book, error) {
	return r.searchByField("author", author)
}

func (r *GORMBookRepository) SearchByTitle(title string) ([]*Book, error) {
	return r.searchByField("title", title)
}

func (r *GORMBookRepository) searchByField(field, value string) ([]*Book, error) {
	var books []*Book
	pattern := "%" + strings.ToLower(value) + "%"
	if err := r.db.Where("LOWER("+field+") LIKE ?", pattern).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
