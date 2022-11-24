package psqlclient

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/jmoiron/sqlx"
	"neatly/pkg/logging"
)

func runUpMigration(db *sqlx.DB, dbname string, migrationsPath string) error {
	logging.GetLogger().Info("Running UP migration...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, dbname, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func RunDownMigrations(db *sqlx.DB, dbname string, migrationsPath string) error {
	logging.GetLogger().Info("Running DOWN migration...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, dbname, driver)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil {
		return err
	}
	return nil
}

func runTestUpMigration(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://db", "neatly", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func RunTestDownMigration(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://db", "neatly", driver)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
