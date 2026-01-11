package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func (c Config) configString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", c.Host, c.Port, c.Database, c.User, c.Password, c.SSLMode)
}

// TODO: Need add login/title indexes because expect often search by those values
// TODO: Add movie_path and movie downloading/streaming options
// TODO: Add auto update total_rating_value, amount_rates
var tableSchema = `
	CREATE TABLE IF NOT EXISTS user_data(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR NOT NULL,
	login VARCHAR NOT NULL UNIQUE,
	password VARCHAR NOT NULL,
	is_admin BOOL NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS movie(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title VARCHAR NOT NULL UNIQUE,
	amount_rates INTEGER NOT NULL DEFAULT 0,
	total_rating_value INTEGER NOT NULL DEFAULT 0,
	rating INTEGER NOT NULL DEFAULT 0, rating REAL GENERATED ALWAYS AS (
		CASE
			WHEN amount_rates = 0 THEN 0
			ELSE total_rating_value / amount_rates
		END
		) STORED,
	created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS favorite_movie(
	user_id UUID REFERENCES user_data ON DELETE CASCADE,
	movie_id UUID REFERENCES movie ON DELETE CASCADE,
	PRIMARY KEY (user_id, movie_id)
	);

	CREATE TABLE IF NOT EXISTS rated_movie(
	user_id UUID REFERENCES user_data ON DELETE CASCADE,
	movie_id UUID REFERENCES movie ON DELETE CASCADE,
	rating INT CHECK (rating between 1 and 10),
	PRIMARY KEY (user_id, movie_id)
	)
	`

// TODO: Delete amount_rates/total_rating_value
// TODO: add materialized view update every N-day at 00:00:00
var functionsSchema = `
	CREATE MATERIALIZED VIEW movie_rating AS
	SELECT AVG(rating)
	FROM rated_movie
	GROUP BY movie_id
`

// WARN: need close() db in outer function
func InitDB(c Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", c.configString())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// TODO: parametrize parameters?
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	_, err = db.Exec(tableSchema)
	if err != nil {
		return nil, fmt.Errorf("create non-exist schema tables: %w", err)
	}
	return db, nil
}
