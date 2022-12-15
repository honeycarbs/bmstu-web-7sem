package testutils

import (
	"neatly/pkg/dbclient"
	"os"
)

var (
	NewAccountQuery = "INSERT INTO users (name, username, email, password_hash) VALUES ('%v', '%v', '%v', '%v') RETURNING id"
)

func Setup(migrationsPath string) (*dbclient.Client, error) {
	client, err := dbclient.NewTestClient(migrationsPath)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Cleanup(client *dbclient.Client, migrationsPath string) error {
	err := client.TestClientClose(migrationsPath)
	if err != nil {
		return err
	}
	return nil
}

func CleanupLogs() error {
	err := os.RemoveAll("build")
	if err != nil {
		return err
	}

	return nil
}
