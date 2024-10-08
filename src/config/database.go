package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const maxRetries = 5

func Setup() *sql.DB{
	var err error
	dsn := "user:password@tcp(mysql:3306)/aws-intern" 

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		log.Printf("Unable to connect to database, retrying... (%d/%d)", i+1, maxRetries)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}

	return db
}
