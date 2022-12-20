package testutils

import (
	"database/sql"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model"
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

func BenchDBSetup(db *sql.DB) {
	createQu := `CREATE TABLE users (
    	id SERIAL NOT NULL UNIQUE,
    	name VARCHAR(255) NOT NULL,
    	username VARCHAR(255) NOT NULL,
    	email VARCHAR(255) NOT NULL,
    	password_hash VARCHAR(255) NOT NULL
	)`

	_, err := db.Exec(createQu)
	if err != nil {
		panic(err)
	}

	entries := []interface{}{}
	for i := 0; i < 10000; i++ {
		entry := createFakeUser(i)
		entries = append(entries, entry)
	}
	insertQu := "INSERT INTO users (id, name, username, email, password_hash) VALUES (:id,:name,:username,:email,:password_hash)"
	dbx := sqlx.NewDb(db, "postgres")
	query, queryArgs, _ := dbx.BindNamed(insertQu, entries)
	query = dbx.Rebind(query)
	_, err = dbx.Queryx(query, queryArgs...)
	if err != nil {
		panic(err)
	}
}

func createFakeUser(i int) model.Account {
	pass := gofakeit.Password(true, true, true, false, false, 6)
	phash, err := model.GeneratePasswordHash(pass)
	if err != nil {
		panic(err)
	}

	return model.Account{
		ID:           i,
		Name:         gofakeit.Name(),
		Username:     gofakeit.Username(),
		Email:        gofakeit.Email(),
		Password:     pass,
		PasswordHash: phash,
	}
}
