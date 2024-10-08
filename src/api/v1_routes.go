package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"aws-intern/handlers"
)

func SetupV1Router(r *gin.Engine, db *sql.DB) {
	v1 := r.Group("/v1")
	{
		v1.POST("/stocks", handlers.UpsertStock(db))
		v1.GET("/stocks", handlers.CheckAllStock(db))
		v1.GET("/stocks/:name", handlers.CheckStock(db))
		v1.DELETE("/stocks", handlers.DeleteAllData(db))

		v1.POST("/sales", handlers.ProcessSale(db))
		v1.GET("/sales", handlers.CalculateSales(db))
	}
}
