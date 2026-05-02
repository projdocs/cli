package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

// _migrate applies all pending up migrations. It is safe to call on every
// startup—migrate tracks applied versions in a schema_migrations table and
// is a no-op when the database is already current.
func _migrate(db *sql.DB, log *zap.Logger) error {
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}
	defer src.Close()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("migrate version: %w", err)
	}

	log.Info("migrations applied",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty),
	)

	return nil
}
