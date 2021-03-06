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
		`CREATE TABLE IF NOT EXISTS server
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
		`CREATE TABLE IF NOT EXISTS game (
			id serial NOT NULL,
			key text NOT NULL,
			mode int,
			player_count int,
			map text,
			mastermode int,
			time_ended timestamp DEFAULT now(),
			PRIMARY KEY (id)
		);`,
		`CREATE INDEX IF NOT EXISTS game_key_time_ended_idx ON game (key, time_ended);`,
		`CREATE TABLE IF NOT EXISTS game_player (
			game_id int NOT NULL,
			name text NOT NULL,
			ip inet NOT NULL,
			country text,
			country_iso text,
			team text,
			kills int,
			deaths int,
			flags int,
			teamkills int,
			accuracy int,
			gun_selected int
		);`,
		`CREATE INDEX IF NOT EXISTS game_player_game_id_idx ON game_player (game_id);`,
		`CREATE INDEX IF NOT EXISTS game_player_name_idx ON game_player (name);`,
		`CREATE INDEX IF NOT EXISTS game_player_ip_idx ON game_player (ip);`,
	}

	for _, t := range tables {
		Conn.MustExec(t)
	}
}
