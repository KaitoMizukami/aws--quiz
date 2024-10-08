package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"aws-intern/api"
	"aws-intern/config"
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

	ipAddress := "54.249.44.164"
	port := "80"
	if err := r.Run(ipAddress + ":" + port); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
