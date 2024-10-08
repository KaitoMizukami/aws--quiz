package models

import (
	"database/sql"
)

func IsExistProduct(db *sql.DB, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE name=?)`
	err := db.QueryRow(query, name).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func CheckStockByName(db *sql.DB, name string) (int, error) {
	var currentQuantity int
	query := `SELECT stockQuantity FROM products WHERE name = ?`
	err := db.QueryRow(query, name).Scan(&currentQuantity)
	if err != nil {
		return -1, err
	}
	return currentQuantity, nil
}

func CheckAllProducStock(db *sql.DB) (map[string]int, error) {
	query := `SELECT name, stockQuantity FROM products`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]int)

	for rows.Next() {
		var name string
		var stock int
		if err := rows.Scan(&name, &stock); err != nil {
			return nil, err
		}
		if stock != 0 {
			data[name] = stock
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func CalculateTotalSales(db *sql.DB) (float64, error) {
	var totalSales float64
	query := `SELECT SUM(totalAmount) FROM salesHistory`
	err := db.QueryRow(query).Scan(&totalSales)
	if err != nil {
		return 0, err
	}
	return totalSales, nil
}

func GetProductID(db *sql.DB, name string) (int, error) {
	var productID int
	query := `SELECT id FROM products WHERE name = ?`
	err := db.QueryRow(query, name).Scan(&productID)
	if err != nil {
		return -1, err
	}
	return productID, nil
}
