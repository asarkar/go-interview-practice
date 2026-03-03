# Learning Materials for RESTful Book Management API

## Building RESTful APIs in Go

This challenge focuses on implementing a RESTful API for managing books, covering core concepts of API design, routing, JSON handling, and database interactions.

### HTTP Server Basics

Go's standard library provides everything needed to build an HTTP server:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Define route handlers
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    
    // Start the server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### RESTful API Design

REST (Representational State Transfer) is an architectural style for designing networked applications:

1. **Resources**: Identified by URLs (e.g., `/books`, `/books/123`)
2. **HTTP Methods**: Used for operations
   - `GET`: Retrieve a resource
   - `POST`: Create a new resource
   - `PUT`: Update an existing resource
   - `DELETE`: Remove a resource
3. **Representations**: Usually JSON or XML
4. **Statelessness**: Each request contains all information needed

### HTTP Routers

While Go's standard library includes basic routing, a router package like `gorilla/mux` offers more flexibility:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    
    // Define routes
    r.HandleFunc("/books", getBooks).Methods("GET")
    r.HandleFunc("/books", createBook).Methods("POST")
    r.HandleFunc("/books/{id}", getBook).Methods("GET")
    r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
    r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
    
    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    
    // Start the server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

### JSON Handling

Go makes it easy to work with JSON using the `encoding/json` package:

```go
// Define a struct for your resource
type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}

// Parsing JSON request body
func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    
    // Decode JSON from request body
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&book); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    
    // Generate a unique ID
    book.ID = uuid.New().String()
    
    // Save the book (implementation depends on your storage)
    books = append(books, book)
    
    // Respond with the created book
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

// Returning JSON response
func getBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}
```

### Route Parameters

Extract parameters from URLs:

```go
func getBook(w http.ResponseWriter, r *http.Request) {
    // Get the ID from the URL
    vars := mux.Vars(r)
    id := vars["id"]
    
    // Find the book
    for _, book := range books {
        if book.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(book)
            return
        }
    }
    
    // Book not found
    http.Error(w, "Book not found", http.StatusNotFound)
}
```

### Query Parameters

Parse query parameters for filtering, pagination, etc.:

```go
func getBooks(w http.ResponseWriter, r *http.Request) {
    // Get query parameters
    query := r.URL.Query()
    
    // Filter by author (if provided)
    author := query.Get("author")
    
    // Parse pagination parameters
    page, _ := strconv.Atoi(query.Get("page"))
    if page < 1 {
        page = 1
    }
    
    limit, _ := strconv.Atoi(query.Get("limit"))
    if limit < 1 || limit > 100 {
        limit = 10 // Default limit
    }
    
    // Apply filters and pagination
    var result []Book
    // ... implementation details ...
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

### Data Storage Options

Several options for storing book data:

#### 1. In-Memory Storage

Simple but non-persistent:

```go
var books []Book // Global variable to store books

// Add a book
books = append(books, book)

// Find a book
for i, book := range books {
    if book.ID == id {
        return book, i, nil
    }
}

// Update a book
books[index] = updatedBook

// Delete a book
books = append(books[:index], books[index+1:]...)
```

#### 2. SQL Database

Using the `database/sql` package with a driver like `go-sql-driver/mysql`:

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/bookstore")
    if err != nil {
        log.Fatal(err)
    }
    
    // Check the connection
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
}

// Create a book
func createBookDB(book Book) (string, error) {
    query := `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`
    result, err := db.Exec(query, book.Title, book.Author, book.Year)
    if err != nil {
        return "", err
    }
    
    id, err := result.LastInsertId()
    return strconv.FormatInt(id, 10), err
}

// Get all books
func getBooksDB() ([]Book, error) {
    query := `SELECT id, title, author, year FROM books`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var books []Book
    for rows.Next() {
        var book Book
        if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
            return nil, err
        }
        books = append(books, book)
    }
    
    return books, nil
}
```

#### 3. NoSQL Database

Using MongoDB with the official Go driver:

```go
import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func initMongoDB() {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    
    collection = client.Database("bookstore").Collection("books")
}

// Create a book
func createBookMongo(book Book) (string, error) {
    book.ID = primitive.NewObjectID().Hex()
    _, err := collection.InsertOne(context.Background(), book)
    return book.ID, err
}

// Get all books
func getBooksMongo() ([]Book, error) {
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    var books []Book
    if err := cursor.All(context.Background(), &books); err != nil {
        return nil, err
    }
    
    return books, nil
}
```

### Middleware

Middleware functions process requests before they reach your handlers:

```go
// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
    })
}

// Authentication middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Validate token...
        
        next.ServeHTTP(w, r)
    })
}

// Apply middleware
r := mux.NewRouter()
r.Use(loggingMiddleware)

// Apply authentication only to certain routes
protected := r.PathPrefix("/api").Subrouter()
protected.Use(authMiddleware)
protected.HandleFunc("/books", createBook).Methods("POST")
```

### Input Validation

Validate incoming data to ensure it meets your requirements:

```go
func validateBook(book Book) error {
    if book.Title == "" {
        return errors.New("title is required")
    }
    
    if book.Author == "" {
        return errors.New("author is required")
    }
    
    if book.Year < 0 || book.Year > time.Now().Year() {
        return errors.New("invalid year")
    }
    
    return nil
}

// Use in your handler
func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    
    if err := validateBook(book); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Continue with book creation...
}
```

### Error Handling

Consistent error handling improves API usability:

```go
// Custom error response
type ErrorResponse struct {
    StatusCode int    `json:"-"`
    Message    string `json:"message"`
    Error      string `json:"error,omitempty"`
}

// Helper function to respond with an error
func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
    response := ErrorResponse{
        StatusCode: statusCode,
        Message:    message,
    }
    
    if err != nil {
        response.Error = err.Error()
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}

// Usage in a handler
func getBook(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    
    book, err := getBookByID(id)
    if err != nil {
        if err == ErrBookNotFound {
            respondWithError(w, http.StatusNotFound, "Book not found", nil)
        } else {
            respondWithError(w, http.StatusInternalServerError, "Failed to get book", err)
        }
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(book)
}
```

### CORS (Cross-Origin Resource Sharing)

Allow requests from different domains:

```go
// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Apply middleware
r.Use(corsMiddleware)
```

### Testing REST APIs

Testing API endpoints with Go's testing package:

```go
func TestGetBooks(t *testing.T) {
    // Create a request
    req, err := http.NewRequest("GET", "/books", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create the handler
    router := mux.NewRouter()
    router.HandleFunc("/books", getBooks).Methods("GET")
    
    // Serve the request
    router.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }
    
    // Check the response body
    var books []Book
    if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
        t.Fatal(err)
    }
    
    // Verify the response
    if len(books) != 2 {
        t.Errorf("expected 2 books, got %d", len(books))
    }
}
```

## Best Practices for RESTful APIs

1. **Use Proper HTTP Status Codes**: 200 for success, 201 for creation, 400 for bad request, 404 for not found, etc.
2. **Consistent Naming Conventions**: Use plural nouns for resources (e.g., `/books` instead of `/book`)
3. **API Versioning**: Include version in the URL or header (e.g., `/api/v1/books`)
4. **Pagination**: Implement pagination for large collections
5. **Filtering, Sorting, and Searching**: Support via query parameters
6. **Documentation**: Use tools like Swagger to document your API

## Further Reading

- [RESTful API Design Guidelines](https://restfulapi.net/)
- [Go Web Examples](https://gowebexamples.com/)
- [Build RESTful APIs with Gorilla Mux](https://www.digitalocean.com/community/tutorials/how-to-make-an-api-with-go-using-gorilla-mux)
- [Go Database Tutorial](https://tutorialedge.net/golang/golang-mysql-tutorial/) 