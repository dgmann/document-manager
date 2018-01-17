package services

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	log "github.com/sirupsen/logrus"
)

func MigratePostgres(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"dialect": "postgres",
			"error":   err,
		}).Fatal("Migration: Failed to connect to database")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	defer m.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"dialect": "postgres",
			"error":   err,
		}).Fatal("Migration: Failed to connect to database")
	}

	m.Up()
	if err != nil {
		log.WithFields(log.Fields{
			"dialect": "postgres",
			"error":   err,
		}).Fatal("Migration failed")
	}
}
