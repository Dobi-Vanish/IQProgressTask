package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(pool *sql.DB) *PostgresRepository {
	db = pool
	return &PostgresRepository{
		Conn: pool,
	}
}

type Transactions struct {
	ID             int             `json:"ID"`
	UserIDSource   int             `json:"UserIDSource"`
	UserIDEndpoint int             `json:"UserIDEndpoint"`
	Amount         decimal.Decimal `json:"Amount"`
	CreatedAt      time.Time       `json:"CreatedAt"`
}

// AddMoney adds some amount of money to users balance
func (u *PostgresRepository) AddMoney(id int, amount decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `UPDATE users SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, stmt, amount, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DecreaseMoney decreasing users balance for some amount
func (u *PostgresRepository) DecreaseMoney(idSource int, amount decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentBalance decimal.Decimal
	err = tx.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = $1", idSource).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	newBalance := currentBalance.Sub(amount)
	if newBalance.IsNegative() {
		return fmt.Errorf("insufficient balance: cannot decrease balance by %s, current balance is %s", amount.String(), currentBalance.String())
	}

	stmt := `UPDATE users SET balance = balance - $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, stmt, amount, time.Now(), idSource)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetLastTransactions returns a slice of 10 transactions, sorted by date
func (u *PostgresRepository) GetLastTransactions(id int) ([]*Transactions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
        SELECT id, useridsource, useridendpoint, amount, createdat
        FROM transactions
        WHERE useridsource = $1 OR useridendpoint = $1
        ORDER BY createdat DESC
        LIMIT 10
    `
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*Transactions
	for rows.Next() {
		var transaction Transactions
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserIDSource,
			&transaction.UserIDEndpoint,
			&transaction.Amount,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transactions, nil
}

// AddTransaction adds transaction to the database
func (u *PostgresRepository) AddTransaction(amount decimal.Decimal, id ...int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt := `
        INSERT INTO transactions (userIDSource, userIDEndpoint, amount, createdat)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var newID int
	err = tx.QueryRowContext(ctx, stmt, id[0], id[1], amount, time.Now()).Scan(&newID)
	if err != nil {
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
