// Package db handles database migrations.
package db

import (
	"embed"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*
var migrationsFS embed.FS

// RunMigrations executes all pending SQL migration files.
func RunMigrations(databaseURL string) {
	// Create an iofs source driver from the embedded filesystem
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create iofs source driver")
	}

	// Pass the iofs source driver to migrate.New
	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Failed to apply database migrations")
	}

	log.Info().Msg("Database migrations applied successfully")
}
