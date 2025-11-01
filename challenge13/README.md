[View the Scoreboard](SCOREBOARD.md)

# Challenge 13: SQL Database Operations with Go

In this challenge, you will implement a product inventory system using Go and SQL. You'll create functions that interact with a SQLite database to perform CRUD operations (Create, Read, Update, Delete) on products.

## Requirements

1. Create a SQLite database with a `products` table  
2. Implement the following functions:  
   - `CreateProduct` - adds a new product to the database  
   - `GetProduct` - retrieves a product by ID  
   - `UpdateProduct` - updates a product's details  
   - `DeleteProduct` - removes a product  
   - `ListProducts` - lists all products with optional filtering  
3. Ensure proper error handling for database operations  
4. Implement transaction support for operations that modify multiple records  
5. Use parameter binding to prevent SQL injection  
6. The included test file has scenarios checking all CRUD operations and error handling 