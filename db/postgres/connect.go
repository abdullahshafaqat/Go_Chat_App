package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	maxRetries    = 5
	retryInterval = 2 * time.Second
)

func Database() *sqlx.DB {
	if os.Getenv("GO_ENV") != "production" {
		log.Println("Loading .env file for local development...")
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (this is OK in production)")
		}
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")


	log.Println("Using DB config:")
	log.Println("DB_HOST:", dbHost)
	log.Println("DB_PORT:", dbPort)
	log.Println("DB_USER:", dbUser)
	log.Println("DB_PASSWORD:", dbPassword)
	log.Println("DB_Name:", dbName)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var db *sqlx.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Successfully connected to PostgreSQL database.")
			return db
		}

		log.Printf(" Database connection attempt %d failed: %v", i, err)
		if i < maxRetries {
			log.Printf("🔁 Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	log.Fatalf(" Could not connect to the database after %d attempts.", maxRetries)
	return nil
}
