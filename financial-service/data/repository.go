package data

import "github.com/shopspring/decimal"

type Repository interface {
	GetLastTransactions(id int) ([]*Transactions, error)
	AddMoney(id int, amount decimal.Decimal) error
	DecreaseMoney(idSource int, amount decimal.Decimal) error
	AddTransaction(amount decimal.Decimal, id ...int) error
}
