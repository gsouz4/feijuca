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

func (s *transactionService) Save(ctx context.Context, clientID int, value int, transactionType string, description string) (entity.Client, error) {
	if clientID < 1 || clientID > 5 {
		return entity.Client{}, errors.New("no rows in result set")
	}
	
	transaction := entity.Transaction{
		Value:       value,
		Type:        transactionType,
		Description: description,
		ClientID:    clientID,
	}

	client, err := s.repo.Save(ctx, transaction)
	if err != nil {
		return entity.Client{}, err
	}

	return client, nil
}

func (s *transactionService) FindBalance(ctx context.Context, clientID int) (entity.BankStatement, error) {
	if clientID < 1 || clientID > 5 {
		return entity.BankStatement{}, errors.New("no rows in result set")
	}

	return s.repo.FindBankStatement(ctx, clientID)
}
