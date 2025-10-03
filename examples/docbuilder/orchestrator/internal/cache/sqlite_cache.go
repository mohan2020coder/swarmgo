package cache

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteCache struct {
	db *sql.DB
}

func NewSQLiteCache(path string) (*SQLiteCache, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	create := `
    CREATE TABLE IF NOT EXISTS kv (
        key TEXT PRIMARY KEY,
        value TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &SQLiteCache{db: db}, nil
}

func (s *SQLiteCache) Get(key string) (string, bool) {
	var value string
	row := s.db.QueryRow("SELECT value FROM kv WHERE key = ?", key)
	if err := row.Scan(&value); err != nil {
		return "", false
	}
	return value, true
}

func (s *SQLiteCache) Set(key, value string) error {
	_, err := s.db.Exec("INSERT OR REPLACE INTO kv(key,value, created_at) VALUES(?,?,?)", key, value, time.Now())
	return err
}
