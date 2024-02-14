package outbounds

import (
	"context"
	"feijuca/domain/entity"
)

type Transaction interface {
	Save(ctx context.Context, transaction entity.Transaction) (entity.Client, error)
	FindBankStatement(ctx context.Context, clientID int) (entity.BankStatement, error)
}
