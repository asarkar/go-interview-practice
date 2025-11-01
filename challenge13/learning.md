# Learning Materials for SQL Database Operations with Go

## SQL Databases in Go

Go provides excellent support for SQL databases through the standard `database/sql` package. This challenge focuses on implementing CRUD operations with SQLite, but the concepts apply to other SQL databases like MySQL, PostgreSQL, etc.

### The database/sql Package

The `database/sql` package provides a generic interface around SQL (or SQL-like) databases. It:

- Manages connection pools 
- Handles transactions
- Provides prepared statements
- Offers drivers for various databases

```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // Note the underscore import
)
```

The underscore import (`_`) is used to import a package solely for its side effects (in this case, registering a database driver).

### Opening a Database Connection

```go
db, err := sql.Open("sqlite3", "path/to/database.db")
if err != nil {
    return nil, err
}

// Test the connection
if err = db.Ping(); err != nil {
    return nil, err
}

return db, nil
```

The `sql.Open()` function doesn't actually establish a connection to the database initially. It's only when you call methods like `Ping()` or execute a query that a connection is established.

### Executing Simple Queries

You can execute simple SQL statements that don't return rows using `db.Exec()`:

```go
result, err := db.Exec(
    "CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)"
)
if err != nil {
    return err
}
```

For `INSERT` operations, you can retrieve the last inserted ID:

```go
result, err := db.Exec(
    "INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
    product.Name, product.Price, product.Quantity, product.Category
)
if err != nil {
    return err
}

// Get the ID of the inserted row
id, err := result.LastInsertId()
if err != nil {
    return err
}
product.ID = id
```

### Querying for Data

To query and retrieve data, use `db.Query()` or `db.QueryRow()`:

```go
// For multiple rows
rows, err := db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
if err != nil {
    return nil, err
}
defer rows.Close() // Always close rows when done

var products []*Product
for rows.Next() {
    p := &Product{}
    err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
    if err != nil {
        return nil, err
    }
    products = append(products, p)
}

// Check for errors from iterating over rows
if err = rows.Err(); err != nil {
    return nil, err
}

return products, nil
```

For a single row, use `QueryRow()`:

```go
row := db.QueryRow("SELECT id, name, price, quantity, category FROM products WHERE id = ?", id)

p := &Product{}
err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("product with ID %d not found", id)
    }
    return nil, err
}

return p, nil
```

### Prepared Statements

For queries that will be executed multiple times, you can use prepared statements to improve performance:

```go
stmt, err := db.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
if err != nil {
    return err
}
defer stmt.Close()

for id, quantity := range updates {
    _, err := stmt.Exec(quantity, id)
    if err != nil {
        return err
    }
}
```

### Transactions

Transactions ensure that a group of operations either all succeed or all fail together:

```go
// Begin a transaction
tx, err := db.Begin()
if err != nil {
    return err
}
defer func() {
    if err != nil {
        tx.Rollback() // Rollback on error
    }
}()

stmt, err := tx.Prepare("UPDATE products SET quantity = ? WHERE id = ?")
if err != nil {
    return err
}
defer stmt.Close()

for id, quantity := range updates {
    result, err := stmt.Exec(quantity, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product with ID %d not found", id)
    }
}

// Commit the transaction
return tx.Commit()
```

### Parameter Binding and SQL Injection Prevention

Always use parameter binding instead of string concatenation to prevent SQL injection:

```go
// DON'T DO THIS - vulnerable to SQL injection
query := fmt.Sprintf("SELECT * FROM products WHERE category = '%s'", category)

// DO THIS - uses parameter binding
rows, err := db.Query("SELECT * FROM products WHERE category = ?", category)
```

Different database drivers use different placeholder styles:

- SQLite, MySQL: `?`
- PostgreSQL: `$1`, `$2`, etc.
- Oracle: `:name`

### Working with NULL Values

SQL databases can contain NULL values. Go provides special types in the `database/sql` package to handle these:

```go
import (
    "database/sql"
)

type Product struct {
    ID       int64
    Name     string
    Price    float64
    Quantity int
    Category sql.NullString // Can be NULL
}

// When scanning
var category sql.NullString
err := row.Scan(&id, &name, &price, &quantity, &category)

// When using
if category.Valid {
    fmt.Println(category.String)
} else {
    fmt.Println("Category is NULL")
}
```

### Error Handling

There are several error types to check for:

```go
if err == sql.ErrNoRows {
    // No rows returned (not necessarily an error)
    return nil, fmt.Errorf("product not found")
}

// Check for unique constraint violation
if strings.Contains(err.Error(), "UNIQUE constraint failed") {
    return nil, fmt.Errorf("product with that name already exists")
}
```

### Connection Pooling

The `database/sql` package automatically handles connection pooling. You can control pool behavior:

```go
db.SetMaxOpenConns(25)  // Maximum number of open connections
db.SetMaxIdleConns(25)  // Maximum number of idle connections
db.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused
```

### Best Practices

1. Always close resources like rows and statements
2. Use transactions for operations that must succeed as a group
3. Never concatenate strings to build SQL queries
4. Check for specific errors like sql.ErrNoRows
5. Keep database connections open for the lifetime of your application
6. Consider using ORMs or query builders for complex applications 