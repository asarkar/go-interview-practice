package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Product represents a product in the inventory system
type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ProductStore manages product operations
type ProductStore struct {
	db *sql.DB
}

// NewProductStore creates a new ProductStore with the given database connection
func NewProductStore(db *sql.DB) *ProductStore {
	if err := resetProductsTable(db); err != nil {
		panic(err)
	}
	return &ProductStore{db: db}
}

func resetProductsTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS products;
		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			price REAL,
			quantity INTEGER,
			category TEXT
		);
	`)
	return err
}

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	return sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=memory", dbPath))
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := ps.db.ExecContext(
		ctx,
		"INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)",
		product.Name, product.Price, product.Quantity, product.Category,
	)
	if err != nil {
		return err
	}

	product.ID, err = res.LastInsertId()
	return err
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	row := ps.db.QueryRowContext(
		ctx,
		"SELECT id, name, price, quantity, category FROM products WHERE id = ?",
		id,
	)

	var p Product
	if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(p *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := ps.db.ExecContext(ctx,
		"UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?",
		p.Name, p.Price, p.Quantity, p.Category, p.ID,
	)
	if err != nil {
		return err
	}

	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numRowsAffected != 1 {
		return fmt.Errorf(
			"expected one row with id %d to be affected, but affected %d",
			p.ID,
			numRowsAffected,
		)
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := ps.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	numRowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numRowsAffected != 1 {
		return fmt.Errorf(
			"expected one row with id %d to be affected, but affected %d",
			id,
			numRowsAffected,
		)
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "SELECT id, name, price, quantity, category FROM products"
	args := []any{}
	if len(category) > 0 {
		query += " WHERE category = ?"
		args = append(args, category)
	}

	rows, err := ps.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows close failed: %v", cerr)
		}
	}()

	var products []*Product
	for rows.Next() {
		var p Product
		if scanErr := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); scanErr != nil {
			return nil, scanErr
		}
		products = append(products, &p)
	}

	return products, rows.Err()
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := ps.db.BeginTx(ctx, nil) // `nil` uses default options
	if err != nil {
		return err
	}

	// Ensure rollback happens if commit is not reached.
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("%v", rbErr)
		}
	}()

	stmt, err := tx.PrepareContext(ctx, "UPDATE products SET quantity = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()

	var totalAffected int64

	for id, qty := range updates {
		res, err := stmt.ExecContext(ctx, qty, id)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		totalAffected += affected
	}

	if totalAffected != int64(len(updates)) {
		return fmt.Errorf("expected %d products to be updated, but updated %d",
			len(updates), totalAffected)
	}

	return tx.Commit()
}
