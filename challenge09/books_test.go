package books

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *httptest.Server {
	repo := NewInMemoryBookRepository()
	service := NewBookService(repo)
	handler := NewBookHandler(service)
	return httptest.NewServer(handler.Router())
}

func TestGetAllBooksEmpty(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/books", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var books []*Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(books) != 0 {
		t.Errorf("Expected empty array; got %d books", len(books))
	}
}

func TestCreateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a new book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", resp.Status)
	}

	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if createdBook.ID == "" {
		t.Error("Expected book to have an ID")
	}
	if createdBook.Title != book.Title {
		t.Errorf("Expected book title %s; got %s", book.Title, createdBook.Title)
	}
}

func TestCreateBookInvalid(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a book with missing required fields
	book := &Book{
		// Title intentionally missing
		Author: "John Doe",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", resp.Status)
	}
}

func TestGetBookByID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First create a book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
	}

	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode created book: %v", err)
	}
	resp.Body.Close()

	// Now get the book by ID
	resp, err = http.Get(fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var retrievedBook Book
	if err := json.NewDecoder(resp.Body).Decode(&retrievedBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if retrievedBook.ID != createdBook.ID {
		t.Errorf("Expected book ID %s; got %s", createdBook.ID, retrievedBook.ID)
	}
	if retrievedBook.Title != book.Title {
		t.Errorf("Expected book title %s; got %s", book.Title, retrievedBook.Title)
	}
}

func TestGetBookByIDNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/books/nonexistent", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestUpdateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First create a book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
	}

	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode created book: %v", err)
	}
	resp.Body.Close()

	// Now update the book
	updatedBook := createdBook
	updatedBook.Description = "Updated description"

	updatedBookJSON, err := json.Marshal(updatedBook)
	if err != nil {
		t.Fatalf("Failed to marshal updated book: %v", err)
	}
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID),
		bytes.NewBuffer(updatedBookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var returnedBook Book
	if err := json.NewDecoder(resp.Body).Decode(&returnedBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if returnedBook.Description != updatedBook.Description {
		t.Errorf(
			"Expected description %s; got %s",
			updatedBook.Description,
			returnedBook.Description,
		)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	book := &Book{
		ID:            "nonexistent",
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/books/nonexistent", server.URL),
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestUpdateBookInvalid(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// PUT with missing required fields (title)
	book := &Book{
		Author:        "John Doe",
		PublishedYear: 2020,
		ISBN:          "978-1234567890",
		Description:   "A book",
	}
	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/books/some-id", server.URL),
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request for missing title; got %v", resp.Status)
	}
}

func TestPartiallyUpdateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a book first
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}
	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
	}
	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode created book: %v", err)
	}
	resp.Body.Close()

	// PATCH with only description
	desc := "Updated via PATCH"
	patch := BookPatch{Description: &desc}
	patchJSON, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Failed to marshal patch: %v", err)
	}
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID),
		bytes.NewBuffer(patchJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PATCH request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
	var returnedBook Book
	if err := json.NewDecoder(resp.Body).Decode(&returnedBook); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if returnedBook.Description != desc {
		t.Errorf("Expected description %q; got %q", desc, returnedBook.Description)
	}
	if returnedBook.Title != createdBook.Title {
		t.Errorf("Expected title unchanged %q; got %q", createdBook.Title, returnedBook.Title)
	}
}

func TestPartiallyUpdateBookNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	desc := "Updated description"
	patch := BookPatch{Description: &desc}
	patchJSON, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Failed to marshal patch: %v", err)
	}
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/books/nonexistent", server.URL),
		bytes.NewBuffer(patchJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PATCH request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestPartiallyUpdateBookInvalid(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a book first
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "Original description",
	}
	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
	}
	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode created book: %v", err)
	}
	resp.Body.Close()

	// PATCH with empty title (invalid)
	emptyTitle := ""
	patch := BookPatch{Title: &emptyTitle}
	patchJSON, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Failed to marshal patch: %v", err)
	}
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID),
		bytes.NewBuffer(patchJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PATCH request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request for empty title; got %v", resp.Status)
	}
}

func TestDeleteBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First create a book
	book := &Book{
		Title:         "The Go Programming Language",
		Author:        "Alan A. A. Donovan and Brian W. Kernighan",
		PublishedYear: 2015,
		ISBN:          "978-0134190440",
		Description:   "The definitive guide to programming in Go",
	}

	bookJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/api/books", server.URL),
		"application/json",
		bytes.NewBuffer(bookJSON),
	)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
	}

	var createdBook Book
	if err := json.NewDecoder(resp.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode created book: %v", err)
	}
	resp.Body.Close()

	// Now delete the book
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	// Verify the book was deleted
	resp, err = http.Get(fmt.Sprintf("%s/api/books/%s", server.URL, createdBook.ID))
	if err != nil {
		t.Fatalf("Failed to verify deletion: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found after deletion; got %v", resp.Status)
	}
	resp.Body.Close()
}

func TestDeleteBookNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/books/nonexistent", server.URL),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found; got %v", resp.Status)
	}
}

func TestSearchBooksBothParams(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?author=Kernighan&title=Go", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"Expected status Bad Request when both author and title provided; got %v",
			resp.Status,
		)
	}
}

func TestSearchBooksByAuthor(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create several books
	books := []*Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Alan A. A. Donovan and Brian W. Kernighan",
			PublishedYear: 2015,
			ISBN:          "978-0134190440",
			Description:   "The definitive guide to programming in Go",
		},
		{
			Title:         "Go in Action",
			Author:        "William Kennedy",
			PublishedYear: 2015,
			ISBN:          "978-1617291784",
			Description:   "An introduction to Go",
		},
		{
			Title:         "The C Programming Language",
			Author:        "Brian W. Kernighan and Dennis Ritchie",
			PublishedYear: 1988,
			ISBN:          "978-0131103627",
			Description:   "The definitive guide to C",
		},
	}

	for _, book := range books {
		bookJSON, err := json.Marshal(book)
		if err != nil {
			t.Fatalf("Failed to marshal book: %v", err)
		}
		resp, err := http.Post(
			fmt.Sprintf("%s/api/books", server.URL),
			"application/json",
			bytes.NewBuffer(bookJSON),
		)
		if err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Search for books by Kernighan
	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?author=Kernighan", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var foundBooks []*Book
	if err := json.NewDecoder(resp.Body).Decode(&foundBooks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(foundBooks) != 2 {
		t.Errorf("Expected 2 books; got %d", len(foundBooks))
	}
}

func TestSearchBooksByTitle(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create several books
	books := []*Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Alan A. A. Donovan and Brian W. Kernighan",
			PublishedYear: 2015,
			ISBN:          "978-0134190440",
			Description:   "The definitive guide to programming in Go",
		},
		{
			Title:         "Go in Action",
			Author:        "William Kennedy",
			PublishedYear: 2015,
			ISBN:          "978-1617291784",
			Description:   "An introduction to Go",
		},
		{
			Title:         "The C Programming Language",
			Author:        "Brian W. Kernighan and Dennis Ritchie",
			PublishedYear: 1988,
			ISBN:          "978-0131103627",
			Description:   "The definitive guide to C",
		},
	}

	for _, book := range books {
		bookJSON, err := json.Marshal(book)
		if err != nil {
			t.Fatalf("Failed to marshal book: %v", err)
		}
		resp, err := http.Post(
			fmt.Sprintf("%s/api/books", server.URL),
			"application/json",
			bytes.NewBuffer(bookJSON),
		)
		if err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Search for Go books
	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?title=Go", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var foundBooks []*Book
	if err := json.NewDecoder(resp.Body).Decode(&foundBooks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(foundBooks) != 2 {
		t.Errorf("Expected 2 books; got %d", len(foundBooks))
	}
}

func TestSearchBooksByTitleNoResults(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create several books
	books := []*Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Alan A. A. Donovan and Brian W. Kernighan",
			PublishedYear: 2015,
			ISBN:          "978-0134190440",
			Description:   "The definitive guide to programming in Go",
		},
		{
			Title:         "Go in Action",
			Author:        "William Kennedy",
			PublishedYear: 2015,
			ISBN:          "978-1617291784",
			Description:   "An introduction to Go",
		},
	}

	for _, book := range books {
		bookJSON, err := json.Marshal(book)
		if err != nil {
			t.Fatalf("Failed to marshal book: %v", err)
		}
		resp, err := http.Post(
			fmt.Sprintf("%s/api/books", server.URL),
			"application/json",
			bytes.NewBuffer(bookJSON),
		)
		if err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 201 Created; got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Search for Python books
	resp, err := http.Get(fmt.Sprintf("%s/api/books/search?title=Python", server.URL))
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var foundBooks []*Book
	if err := json.NewDecoder(resp.Body).Decode(&foundBooks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(foundBooks) != 0 {
		t.Errorf("Expected 0 books; got %d", len(foundBooks))
	}
}
