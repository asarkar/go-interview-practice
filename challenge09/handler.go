package books

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
	Service BookService
}

// NewBookHandler creates a new book handler
func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{Service: service}
}

// Router returns a chi router with all book endpoints registered
func (h *BookHandler) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/api/books", h.getAllBooks)
	r.Post("/api/books", h.createBook)
	r.Get("/api/books/search", h.searchBooks)
	r.Get("/api/books/{id}", h.getBookByID)
	r.Put("/api/books/{id}", h.updateBook)
	r.Patch("/api/books/{id}", h.partiallyUpdateBook)
	r.Delete("/api/books/{id}", h.deleteBook)
	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, errNotFound) {
		writeError(w, http.StatusNotFound, "book not found")
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func (h *BookHandler) getAllBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := h.Service.GetAllBooks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if books == nil {
		books = []*Book{}
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.Service.CreateBook(&book); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, &book)
}

func (h *BookHandler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := h.Service.GetBookByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	id := chi.URLParam(r, "id")
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.Service.UpdateBook(id, &book); err != nil {
		writeServiceError(w, err)
		return
	}
	book.ID = id
	writeJSON(w, http.StatusOK, &book)
}

func (h *BookHandler) partiallyUpdateBook(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	id := chi.URLParam(r, "id")
	var patch BookPatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	book, err := h.Service.PartiallyUpdateBook(id, &patch)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.Service.DeleteBook(id); err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "book deleted"})
}

func (h *BookHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	author := q.Get("author")
	title := q.Get("title")

	if author != "" && title != "" {
		writeError(w, http.StatusBadRequest, "provide either author or title, not both")
		return
	}

	var books []*Book
	var err error

	switch {
	case author != "":
		books, err = h.Service.SearchBooksByAuthor(author)
	case title != "":
		books, err = h.Service.SearchBooksByTitle(title)
	default:
		books, err = h.Service.GetAllBooks()
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if books == nil {
		books = []*Book{}
	}
	writeJSON(w, http.StatusOK, books)
}
