package outbounds

import (
	"context"
	"feijuca/domain/entity"
)

type Transaction interface {
	Save(ctx context.Context, transaction entity.Transaction) error
	FindBankStatement(ctx context.Context, clientID int) (entity.BankStatement, error)
	CanDebitValue(ctx context.Context, clientID int, value int) (bool, error)
}
