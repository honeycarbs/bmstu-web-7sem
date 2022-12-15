package dbclient

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ory/dockertest"
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

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping DB: %v", err)
	}

	err = runUpMigration(db, cfg.DBName, cfg.MigrationsPath)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}

func NewBenchClient(migrationsPath string) (*Client, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"localhost",
		"5432",
		"testdb",
		"testdb",
		"testdb",
		"disable",
	))
	if err != nil {
		logging.GetLogger().Info("Error while connecting to db")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping DB: %v", err)
	}

	err = runUpMigration(db, "testdb", migrationsPath)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}

func (c *Client) BenchClientClose(migrationsPath string) error {
	err := RunDownMigrations(c.DB, "testdb", migrationsPath)
	if err != nil {
		return err
	}

	err = c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close(cfg session.DB) error {
	err := RunDownMigrations(c.DB, cfg.DBName, cfg.MigrationsPath)
	if err != nil {
		return err
	}

	err = c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewTestClient() (*Client, error) {
	logging.Init()
	db, err := sqlx.Open("sqlite3", "file:db/sqlite.db")
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

func NewIntegrationClinent(resource *dockertest.Resource, migrationPath string) (*Client, error) {
	dsn := "postgres://test:pass@0.0.0.0:" + resource.GetPort("5432/tcp") + "/test?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = runUpMigration(db, "test", migrationPath)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
