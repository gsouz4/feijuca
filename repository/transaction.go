package repository

import (
	"context"
	"database/sql"
	"log"

	"feijuca/domain/entity"
	"feijuca/domain/ports/outbounds"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) outbounds.Transaction {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Save(ctx context.Context, transaction entity.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO trasactions ("value", "type", "description", "client_id") VALUES ($1, $2, $3, $4)`

	_, err = tx.ExecContext(
		ctx,
		query,
		transaction.Value,
		transaction.Type,
		transaction.Description,
		transaction.ClientID,
	)

	if err != nil {
		return err
	}

	query = `UPDATE clients SET "balance" = $1 + clients.balance WHERE id = $2;`

	if transaction.Type == "d" {
		query = `UPDATE clients SET "balance" = $1 - clients.balance WHERE id = $2;`
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		transaction.Value,
		transaction.ClientID,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
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

	rows, err := r.db.QueryContext(ctx, query, clientID)
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

		transactions = append(transactions, transaction)
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
	query := `SELECT "balance", "limit" FROM clients WHERE id = $1`

	var balance entity.Balance

	err := r.db.QueryRowContext(ctx, query, clientID).Scan(&balance.Total, &balance.Limit)

	return balance, err
}

func (r *transactionRepository) CanDebitValue(ctx context.Context, clientID int, value int) (canDebit bool, err error) {
	query := `SELECT "balance" - $1 < -"limit" FROM clients WHERE client_id = $2 FOR UPDATE;`

	err = r.db.QueryRowContext(ctx, query, value, clientID).Scan(&canDebit)

	return canDebit, err
}
