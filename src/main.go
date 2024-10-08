package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"aws-intern/config"
	"aws-intern/api"
)

func main() {
	db := config.Setup()
	defer db.Close()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "AWS")
	})
	r.GET("/secret", func(c *gin.Context) {
		c.String(200, "SUCCESS")
	})

	api.SetupV1Router(r, db)

	r.Run(":8080")
}
