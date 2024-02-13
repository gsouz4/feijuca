package inbounds

import (
	"context"
	"feijuca/domain/entity"
)

type TransactionService interface {
	Save(ctx context.Context, clientID int, value int, transactionType string, description string) error
	FindBalance(ctx context.Context, clientID int) (entity.BankStatement, error)
}
