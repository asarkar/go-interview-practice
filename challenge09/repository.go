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
	result := r.db.Model(&Book{}).Where("id = ?", id).Updates(map[string]any{
		"title":          book.Title,
		"author":         book.Author,
		"published_year": book.PublishedYear,
		"isbn":           book.ISBN,
		"description":    book.Description,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errNotFound
	}
	return nil
}

func (r *GORMBookRepository) Patch(id string, patch *BookPatch) (*Book, error) {
	updates := make(map[string]any)
	if patch.Title != nil {
		updates["title"] = *patch.Title
	}
	if patch.Author != nil {
		updates["author"] = *patch.Author
	}
	if patch.PublishedYear != nil {
		updates["published_year"] = *patch.PublishedYear
	}
	if patch.ISBN != nil {
		updates["isbn"] = *patch.ISBN
	}
	if patch.Description != nil {
		updates["description"] = *patch.Description
	}
	if len(updates) == 0 {
		return r.GetByID(id)
	}
	result := r.db.Model(&Book{}).Where("id = ?", id).Updates(updates)
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
