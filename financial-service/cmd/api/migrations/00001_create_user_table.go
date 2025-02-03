package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUserTable, downCreateUserTable)
}

func upCreateUserTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec("CREATE TABLE IF NOT EXISTS users (\n    id SERIAL PRIMARY KEY,\n    email VARCHAR(255) NOT NULL UNIQUE,\n    first_name VARCHAR(100),\n    last_name VARCHAR(100),\n    password VARCHAR(255) NOT NULL,\n    active INT NOT NULL DEFAULT 1,\n    balance DECIMAL(15, 2) NOT NULL DEFAULT 0, \n    created_at TIMESTAMP NOT NULL DEFAULT NOW(), \n    updated_at TIMESTAMP NOT NULL DEFAULT NOW() \n);")
	if err != nil {
		return err
	}
	return nil
}

func downCreateUserTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		return err
	}
	return nil
}
