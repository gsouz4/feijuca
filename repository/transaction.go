package repository

import (
	"context"
	"database/sql"

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

func (*transactionRepository) Save(ctx context.Context, transaction entity.Transaction) error {

	return nil
}

func (*transactionRepository) FindBankStatement(ctx context.Context, clientID int) {

}