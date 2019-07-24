package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	// Import for side effects
	_ "github.com/lib/pq"
)

var (
	// Conn is the database connection
	Conn *sqlx.DB
)

// Setup initializes the DB connection
func Setup() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	Conn = sqlx.MustConnect("postgres", dbConnectionString)
	createTables()
}

func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS servers
		(
			key text UNIQUE,
			ip inet,
			port int,
			description text,
			max_clients int,
			country text,
			country_iso text,
			match bool,
			custom_connect text,
			on_ms bool,
			disabled bool,
			dead bool,
			PRIMARY KEY (key)
		)`,
	}

	for _, t := range tables {
		Conn.MustExec(t)
	}
}
