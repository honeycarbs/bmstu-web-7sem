package psqlclient

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"neatly/internal/session"
	"neatly/pkg/logging"
)

type Client struct {
	DB *sqlx.DB
}

func NewClient(cfg session.DB) (*Client, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	))
	if err != nil {
		logging.GetLogger().Info("Error while connecting to db")
		return nil, err
	}

	err = runUpMigration(db, cfg.DBName, cfg.MigrationsPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}

func NewTestClient() (*Client, error) {
	logging.Init()
	db, err := sqlx.Open("sqlite3", "file:testdb/test.db?parseTime=true")
	if err != nil {
		return nil, err
	}

	err = runTestUpMigration(db)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}

func TestClientClose(client *Client) error {
	err := RunTestDownMigration(client.DB)
	if err != nil {
		return err
	}

	err = client.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func runUpMigration(db *sqlx.DB, dbname string, migrationsPath string) error {
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

func runTestUpMigration(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://testdb", "neatly", driver)
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
	m, err := migrate.NewWithDatabaseInstance("file://testdb", "neatly", driver)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
