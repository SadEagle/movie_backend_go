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

var schema = `
	CREATE TABLE IF NOT EXISTS user_data(
	id UUID PRIMARY KEY,
	name VARCHAR NOT NULL,
	login VARCHAR NOT NULL UNIQUE,
	password VARCHAR NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS movie(
	id UUID PRIMARY KEY,
	title VARCHAR NOT NULL UNIQUE,
	amount_marks INTEGER NOT NULL DEFAULT 0,
	total_mark INTEGER NOT NULL DEFAULT 0,
	rating REAL GENERATED ALWAYS AS (
		CASE
			WHEN amount_marks = 0 THEN -1
			ELSE total_mark / amount_marks
		END
		) STORED,
	created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS favorite_movie(
	user_id UUID REFERENCES user_data ON DELETE CASCADE,
	movie_id UUID REFERENCES movie ON DELETE CASCADE,
	PRIMARY KEY (user_id, movie_id)
	);

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

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("create non-exist schema tables: %w", err)
	}
	return db, nil
}
