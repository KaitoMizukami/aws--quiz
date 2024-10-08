package models

import (
	"database/sql"
)

func CreateProduct(db *sql.DB, name string, quantity int) error {
	query := `INSERT INTO products (name, stockQuantity) VALUES (?, ?)`
	_, err := db.Exec(query, name, quantity)
	return err
}

func AddStock(db *sql.DB, name string, quantity int) error {
	query := `UPDATE products SET stockQuantity = stockQuantity + ? WHERE name = ?`
	_, err := db.Exec(query, quantity, name)
	return err
}

func RemoveStock(db *sql.DB, name string, quantity int) error {
	query := `UPDATE products SET stockQuantity = stockQuantity - ? WHERE name = ?`
	_, err := db.Exec(query, quantity, name)
	return err
}

func CreateStockHistory(db *sql.DB, productID int, quantityChanged int) error {
	query := `INSERT INTO stockHistory (productID, quantityChanged) VALUES (?, ?)`
	_, err := db.Exec(query, productID, quantityChanged)
	return err
}

func CreateSalesHistory(db *sql.DB, productID int, quantitySold int, price float64) error {
	query := `INSERT INTO salesHistory (productID, quantitySold, totalAmount) VALUES (?, ?, ?)`
	_, err := db.Exec(query, productID, quantitySold, float64(quantitySold)*price)
	return err
}

func DeleteAllData(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM stockHistory"); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM salesHistory"); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM products"); err != nil {
		return err
	}
	return nil
}