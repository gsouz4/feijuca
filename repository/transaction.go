package repository

import (
	"context"
	"errors"
	"log"

	"feijuca/domain/entity"
	"feijuca/domain/ports/outbounds"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepository struct {
	db          *pgxpool.Pool
	arruda      string
	portudo     string
	transaction bool
}

func NewTransactionRepository(db *pgxpool.Pool) outbounds.Transaction {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Save(ctx context.Context, transaction entity.Transaction) (client entity.Client, err error) {
	// tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	// if err != nil {
	// 	return client, err
	// }
	// defer tx.Rollback(ctx)

	if transaction.Type == "d" {
		query := `UPDATE clients SET "balance" = clients.balance - $1
			WHERE id = $2
			AND abs("balance" - $3) <= "limit"
			RETURNING "balance", "limit";`

		var limit, balance int

		err := r.db.QueryRow(
			ctx,
			query,
			transaction.Value,
			transaction.ClientID,
			transaction.Value,
		).Scan(&limit, &balance)

		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Client{}, errors.New("invalid transaction")
		}

		if err != nil {
			return client, err
		}
	}

	if transaction.Type == "c" {
		query := `UPDATE clients SET "balance" = clients.balance + $1 WHERE id = $2 RETURNING clients.limit, clients.balance;`

		var limit, balance int

		err := r.db.QueryRow(ctx, query, transaction.Value, transaction.ClientID).Scan(&limit, &balance)

		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Client{}, errors.New("invalid transaction")
		}

		if err != nil {
			return entity.Client{}, err
		}
	}

	query := `INSERT INTO transactions ("value", "type", "description", "client_id") VALUES ($1, $2, $3, $4)`

	_, err = r.db.Exec(ctx, query, transaction.Value, transaction.Type, transaction.Description, transaction.ClientID)
	if err != nil {
		return client, err
	}

	// err = tx.Commit(ctx)
	// if err != nil {
	// 	return client, err
	// }

	return client, nil
}

func (r *transactionRepository) FindBankStatement(ctx context.Context, clientID int) (entity.BankStatement, error) {
	query := `
		SELECT
			"value",
			"type",
			"description",
			"created_at"
		FROM transactions
			WHERE client_id = $1
			ORDER BY "created_at" DESC LIMIT 10`

	transactions := make([]entity.Transaction, 0)

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return entity.BankStatement{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction entity.Transaction

		if err := rows.Scan(
			&transaction.Value,
			&transaction.Type,
			&transaction.Description,
			&transaction.Date,
		); err != nil {
			log.Fatal("aa")
		}

		if transaction.Value > 0 {
			transactions = append(transactions, transaction)
		}
	}

	balance, err := r.FindBalance(ctx, clientID)
	if err != nil {
		return entity.BankStatement{}, err
	}

	var statement entity.BankStatement

	statement.Balance = balance
	statement.LastTransactions = transactions

	return statement, nil
}

func (r *transactionRepository) FindBalance(ctx context.Context, clientID int) (entity.Balance, error) {
	query := `SELECT "balance", "limit" FROM clients WHERE id = $1 FOR UPDATE`

	var balance entity.Balance

	err := r.db.QueryRow(ctx, query, clientID).Scan(&balance.Total, &balance.Limit)

	return balance, err
}
