// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

func main() {
	// Koneksi ke database PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	// Retry connection with backoff
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error opening database: %v. Retrying in %d seconds...", err, (i+1)*2)
			time.Sleep(time.Duration((i+1)*2) * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			log.Printf("Error connecting to database: %v. Retrying in %d seconds...", err, (i+1)*2)
			err := db.Close()
			if err != nil {
				return
			} // Close existing connection before retrying
			time.Sleep(time.Duration((i+1)*2) * time.Second)
			continue
		}
		log.Println("Successfully connected to the database!")
		break
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after multiple retries: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", healthCheck)
	e.GET("/pingdb", pingDB)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9688"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// Handler untuk health check sederhana
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from Go Echo App!")
}

// Handler untuk menguji koneksi database
func pingDB(c echo.Context) error {
	err := db.Ping()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to ping database: %v", err))
	}
	return c.String(http.StatusOK, "Successfully pinged database!")
}
