package psqlclient

import (
	"fmt"
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
	db, err := sqlx.Open("sqlite3", "file:db/test.db")
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
