package handlers

import (
	"fmt"
	"log"
	"math"
	"database/sql"
	"net/http"
	"bytes"
	"strconv"

	"github.com/gin-gonic/gin"

	"aws-intern/models"
)

type SalesRequest struct {
	Name   string   `json:"name" binding:"required"`
	Amount *int     `json:"amount,omitempty"`
	Price  *float64 `json:"price,omitempty"`
}

func ProcessSale(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
		var req SalesRequest

		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Printf("failed to bind json: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "ERROR"})
			return
		}
		amount := 1
		if req.Amount != nil {
			amount = *req.Amount
		}
		if req.Amount != nil && *req.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ERROR"})
			return
		}
		if req.Price != nil && *req.Price < 0.0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ERROR"})
			return
		}

		currentStock, err := models.CheckStockByName(db, req.Name)
		if err != nil {
			log.Printf("failed to check stock by name: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if currentStock < amount {
			log.Printf("Insufficient stock for your desired quantity.")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient stock for your desired quantity."})
			return
		}

		err = models.RemoveStock(db, req.Name, amount)
		if err != nil {
			log.Printf("failed to remove stock: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		productID, err := models.GetProductID(db, req.Name)
		if err != nil {
			log.Printf("failed to get product id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var price float64
		if req.Price != nil {
			price = *req.Price
		}

		err = models.CreateSalesHistory(db, productID, amount, price)
		if err != nil {
			log.Printf("failed to create sales history: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		err = models.CreateStockHistory(db, productID, amount)
		if err != nil {
			log.Printf("failed to create stock history: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		locationURL := fmt.Sprintf("http://localhost:8080/v1/sales/%s", req.Name)
		c.Header("Location", locationURL)

		var buffer bytes.Buffer
		buffer.WriteString("{")
		buffer.WriteString(`"name":"` + req.Name + `"`)
		if req.Amount != nil {
			buffer.WriteString(`,"amount":` + strconv.Itoa(*req.Amount))
		}
		if req.Price != nil {
			buffer.WriteString(`,"price":` + strconv.FormatFloat(*req.Price, 'f', -1, 64))
		}
		buffer.WriteString("}")
		c.Data(http.StatusOK, "application/json", buffer.Bytes())
    }
}

func CalculateSales(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
		totalSales, err := models.CalculateTotalSales(db)
		if err != nil {
			log.Printf("failed to calculate total sales: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		totalSales = math.Ceil(totalSales*100) / 100

		c.JSON(http.StatusOK, gin.H{"sales": fmt.Sprintf("%.2f", totalSales)})
	}
}