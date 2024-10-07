package database

import (
	"errors"
	"time"

	"example_project/assets"
	"example_project/internal/leveledlog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

var logger = leveledlog.Logger.WithID("db-logger")

func New(dsn string, autoMigrate bool) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect to database, dsn: %s, details: %v", dsn, err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if autoMigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}

	return &DB{db}, nil
}
