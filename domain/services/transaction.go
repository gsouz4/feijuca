package services

import (
	"context"
	"errors"
	"feijuca/domain/entity"
	"feijuca/domain/ports/inbounds"
	"feijuca/domain/ports/outbounds"
)

type transactionService struct {
	repo outbounds.Transaction
}

func NewTransactionRepository(repo outbounds.Transaction) inbounds.TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) Save(ctx context.Context, value int, clientID int, transactionType string, description string) error {
	if transactionType == "d" {
		canDebit, err := s.repo.CanDebitValue(ctx, clientID, value)
		if err != nil {
			return err
		}

		if !canDebit {
			return errors.New("invalid transaction")
		}
	}

	transaction := entity.Transaction{
		Value:       value,
		Type:        transactionType,
		Description: description,
		ClientID:    clientID,
	}

	err := s.repo.Save(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *transactionService) FindBalance(ctx context.Context, clientID int) (entity.BankStatement, error) {
	return s.repo.FindBankStatement(ctx, clientID)
}
