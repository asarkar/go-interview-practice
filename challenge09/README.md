# Challenge 9: RESTful Book Management API

## Problem Statement

Implement a RESTful API for a book management system using Go. The API should allow users to perform CRUD operations on books, with data persistence using an in-memory database. This challenge tests your ability to design and implement a complete web service, handle HTTP requests and responses, and manage data persistence.

## Requirements

1. Implement a RESTful API with the following endpoints:
   - `GET /api/books`: Get all books
   - `GET /api/books/{id}`: Get a specific book by ID
   - `POST /api/books`: Create a new book
   - `PUT /api/books/{id}`: Update an existing book
   - `DELETE /api/books/{id}`: Delete a book
   - `GET /api/books/search?author={author}`: Search books by author
   - `GET /api/books/search?title={title}`: Search books by title

2. Implement a `Book` struct with the following fields:
   - `ID`: Unique identifier for the book
   - `Title`: Title of the book
   - `Author`: Author of the book
   - `PublishedYear`: Year the book was published
   - `ISBN`: International Standard Book Number
   - `Description`: Brief description of the book

3. Implement an in-memory database (using Go data structures) to store books.

4. Implement proper error handling and status codes:
   - 200 OK: Successful GET, PUT, DELETE
   - 201 Created: Successful POST
   - 400 Bad Request: Invalid input
   - 404 Not Found: Resource not found
   - 500 Internal Server Error: Server-side error

5. Implement input validation for all endpoints.

6. The API should return responses in JSON format.

## Function Signatures and Interfaces

```go
// Book represents a book in the database
type Book struct {
    ID            string `json:"id"`
    Title         string `json:"title"`
    Author        string `json:"author"`
    PublishedYear int    `json:"published_year"`
    ISBN          string `json:"isbn"`
    Description   string `json:"description"`
}

// BookRepository defines the operations for book data access
type BookRepository interface {
    GetAll() ([]*Book, error)
    GetByID(id string) (*Book, error)
    Create(book *Book) error
    Update(id string, book *Book) error
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
    DeleteBook(id string) error
    SearchBooksByAuthor(author string) ([]*Book, error)
    SearchBooksByTitle(title string) ([]*Book, error)
}

// BookHandler handles HTTP requests for book operations
type BookHandler struct {
    Service BookService
}

// Implement appropriate methods on BookHandler to handle HTTP requests
```

## Project Structure

Your solution should follow a clean architecture with separation of concerns:

```
challenge-9/
├── submissions/
│   └── yourusername/
│       └── solution-template.go
├── api/
│   ├── handlers/
│   │   └── book_handler.go
│   └── middleware/
│       └── logger.go
├── domain/
│   └── models/
│       └── book.go
├── repository/
│   └── book_repository.go
├── service/
│   └── book_service.go
└── main.go
```

## Test Cases

Your solution should handle the following test scenarios:

1. Get all books when the database is empty
2. Create a new book with valid data
3. Create a new book with invalid data (missing required fields)
4. Get a specific book by ID when it exists
5. Get a specific book by ID when it doesn't exist
6. Update a book with valid data
7. Update a book that doesn't exist
8. Delete a book that exists
9. Delete a book that doesn't exist
10. Search for books by author with results
11. Search for books by title with no results

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-9/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required components.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-9/` directory:

```bash
go test -v
```

## Bonus Challenges

For those seeking extra challenges:

1. Add authentication and authorization using JWT
2. Implement pagination for the `GET /api/books` endpoint
3. Add rate limiting middleware
4. Implement filtering and sorting options for book queries
5. Add swagger documentation for the API 