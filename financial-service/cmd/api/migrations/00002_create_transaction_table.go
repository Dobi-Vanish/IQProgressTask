package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTransactionTable, downCreateTransactionTable)
}

func upCreateTransactionTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec("CREATE TABLE IF NOT EXISTS transactions (\n    ID SERIAL PRIMARY KEY,\n    userIDSource INT,\n    userIDEndpoint INT,\n    amount DECIMAL(10, 2) NOT NULL,\n    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n    FOREIGN KEY (userIDSource) REFERENCES Users(ID),\n    FOREIGN KEY (userIDEndpoint) REFERENCES Users(ID)\n);\n")
	if err != nil {
		return err
	}
	return nil
}

func downCreateTransactionTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("DROP TABLE IF EXISTS transactions;")
	if err != nil {
		return err
	}
	return nil
}
