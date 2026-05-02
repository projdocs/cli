package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	_ "modernc.org/sqlite" // registers the "sqlite" driver
)

// Config holds connection parameters for the SQLite database.
type Config struct {
	// Path is the file path to the SQLite database file.
	// Use ":memory:" for an in-memory database.
	Path string

	// MaxOpenConns should be kept at 1 for SQLite to avoid
	// SQLITE_BUSY errors—SQLite's writer lock is process-wide.
	MaxOpenConns int

	// ConnMaxLifetime is the maximum time a connection may be reused.
	ConnMaxLifetime time.Duration
}

// DefaultConfig returns a sensible Config for a file-backed database.
func DefaultConfig(path string) Config {
	return Config{
		Path:            path,
		MaxOpenConns:    1,
		ConnMaxLifetime: time.Hour,
	}
}

// Get opens and validates a SQLite database connection, applying
// recommended PRAGMAs for reliability and performance.
func Get(cfg Config, log *zap.Logger) (*sql.DB, error) {
	// The DSN supports SQLite URI parameters.
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000", cfg.Path)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// SQLite does not benefit from a pool—one writer at a time.
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := ping(db, log); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := _migrate(db, log); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrations: %w", err)
	}

	log.Debug("sqlite connected", zap.String("path", cfg.Path))
	return db, nil
}

// ping verifies the connection is live and logs the SQLite version.
func ping(db *sql.DB, log *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping sqlite: %w", err)
	}

	var version string
	if err := db.QueryRowContext(ctx, "SELECT sqlite_version()").Scan(&version); err != nil {
		return fmt.Errorf("query sqlite version: %w", err)
	}

	log.Debug("sqlite version", zap.String("version", version))
	return nil
}
