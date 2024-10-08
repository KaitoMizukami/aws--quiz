package handlers

import (
	"fmt"
	"log"
	"regexp"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"aws-intern/models"
)

type UpsertStockRequest struct {
	Name   string `json:"name" binding:"required"`
	Amount int    `json:"amount"`
}

func UpsertStock(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
		var req UpsertStockRequest
		isAmountProvided := true

		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Printf("failed to bind json: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "ERROR"})
			return
		}
		if req.Amount == 0 {
			req.Amount = 1
			isAmountProvided = false
		} else if req.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ERROR"})
			return
		}
		if !isValidName(req.Name) {
			log.Printf("invalid name format for input name")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid name format",
				"details": "The name must consist only of uppercase and lowercase alphabets and be up to 8 characters long.",
				"inputReceived": req.Name,
			})
			return
		}

		isExist, err := models.IsExistProduct(db, req.Name)
		if err != nil {
			log.Printf("failed to check the product is exist: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	
		if isExist {
			err = models.AddStock(db, req.Name, req.Amount)
			if err != nil {
				log.Printf("failed to add stock: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		} else {
			err = models.CreateProduct(db, req.Name, req.Amount)
			if err != nil {
				log.Printf("failed to create product: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		productID, err := models.GetProductID(db, req.Name)
		if err != nil {
			log.Printf("failed to get product id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		err = models.CreateStockHistory(db, productID, req.Amount)
		if err != nil {
			log.Printf("failed to create stock history: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		locationURL := fmt.Sprintf("http://localhost:8080/v1/stocks/%s", req.Name)
		c.Header("Location", locationURL)

		if isAmountProvided {
			c.JSON(http.StatusOK, req)
		} else {
			c.JSON(http.StatusOK, gin.H{"name": req.Name})
		}
    }
}

func CheckAllStock(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
		stockQuantities, err := models.CheckAllProducStock(db)
		if err != nil {
			log.Printf("failed to check all stock: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, stockQuantities)
	}
}

func CheckStock(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
		name := c.Param("name")

		if !isValidName(name) {
			log.Printf("invalid name format for input name")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid name format",
				"details": "The name must consist only of uppercase and lowercase alphabets and be up to 8 characters long.",
				"inputReceived": name,
			})
			return
		}

		isExist, err := models.IsExistProduct(db, name)
		if err != nil {
			log.Printf("failed to check the product is exist: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	
		if isExist {
			stockQuantity, err := models.CheckStockByName(db, name)
			if err != nil {
				log.Printf("failed to check stock by name: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{name: stockQuantity})
		} else {
			log.Printf("No product found with the name %s", name)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product Not Found",
				"details": fmt.Sprintf("No product found with the name '%s'. Please check the product name and try again.", name),
				"inputReceived": name,
			})
		}
	}
}

func DeleteAllData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := models.DeleteAllData(db) 
		if err != nil {
			log.Printf("failed to delete all data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
}

func isValidName(name string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z]{1,8}$", name)
	return matched
}