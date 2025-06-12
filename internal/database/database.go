package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDatabase(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return database, nil
}

func (db *DB) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			comment_id INTEGER,
			is_like BOOLEAN NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
			FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
			UNIQUE(user_id, post_id),
			UNIQUE(user_id, comment_id),
			CHECK ((post_id IS NOT NULL AND comment_id IS NULL) OR (post_id IS NULL AND comment_id IS NOT NULL))
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	// Insert default categories
	defaultCategories := []string{"General", "Technology", "Sports", "Entertainment", "News", "Gaming"}
	for _, category := range defaultCategories {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			log.Printf("Failed to insert category %s: %v", category, err)
		}
	}

	return nil
}