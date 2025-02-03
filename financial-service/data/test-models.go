package data

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}

func (u *PostgresTestRepository) AddMoney(id int, amount decimal.Decimal) error {
	return nil
}

func (u *PostgresTestRepository) DecreaseMoney(idSource int, amount decimal.Decimal) error {
	return nil
}

func (u *PostgresTestRepository) GetLastTransactions(id int) ([]*Transactions, error) {
	return nil, nil
}

func (u *PostgresTestRepository) AddTransaction(amount decimal.Decimal, id ...int) error {
	return nil
}
