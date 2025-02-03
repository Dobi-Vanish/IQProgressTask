package data

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"log"
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

// User is the structure which holds one user from the database.
type User struct {
	ID        int             `json:"id"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name,omitempty"`
	LastName  string          `json:"last_name,omitempty"`
	Password  string          `json:"-"`
	Active    int             `json:"active"`
	Balance   decimal.Decimal `json:"score"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Transactions struct {
	ID             int             `json:"ID"`
	userIDSource   int             `json:"UserIDSource"`
	userIDEndpoint int             `json:"UserIDEndpoint"`
	amount         decimal.Decimal `json:"Amount"`
	CreatedAt      time.Time       `json:"CreatedAt"`
}

// AddMoney adds  some points
func (u *PostgresRepository) AddMoney(id int, amount decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
        amount = amount + $1,
        updated_at = $2
		where id = $3
	`

	_, err := db.ExecContext(ctx, stmt,
		amount,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	return nil

}

// DecreaseMoney adds  some points
func (u *PostgresRepository) DecreaseMoney(idSource int, amount decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE users SET 
                 amount = amount - $1, 
    			 updated_at = $2 
             WHERE idEndpoint = $3`

	_, err := db.ExecContext(ctx, stmt,
		amount,
		time.Now(),
		idSource,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresRepository) GetLastTransactions(id int) ([]*Transactions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select ID, userIDSource, userIDEndpoint, amount, CreatedAt
	from transactions where userIDSource = $1 order by CreatedAt desc LIMIT 10`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transactions

	for rows.Next() {
		var transaction Transactions
		err := rows.Scan(
			&transaction.ID,
			&transaction.userIDSource,
			&transaction.userIDEndpoint,
			&transaction.amount,
			&transaction.CreatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (u *PostgresRepository) AddTransaction(amount decimal.Decimal, id ...int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into transactions (userIDSource, userIdEndpoint, amount, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning id`

	err := db.QueryRowContext(ctx, stmt,
		id[1],
		id[2],
		amount,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return err
	}

	return nil
}
