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
	ch   chan TransactionEvent
}

func NewTransactionService(repo outbounds.Transaction, ch chan TransactionEvent) inbounds.TransactionService {
	return &transactionService{
		repo: repo,
		ch:   ch,
	}
}

func (s *transactionService) Save(ctx context.Context, clientID int, value int, transactionType string, description string) (entity.Client, error) {
	transaction := entity.Transaction{
		Value:       value,
		Type:        transactionType,
		Description: description,
		ClientID:    clientID,
	}

	balance, err := s.repo.FindBalance(ctx, clientID)
	if err != nil {
		return entity.Client{}, err
	}

	if transaction.Type == "d" {
		if (balance.Total - value) < balance.Limit {
			return entity.Client{}, errors.New("invalid request")
		}

	}

	var eventType string
	var balanceFinal int

	if transaction.Type == "d" {
		eventType = "Debit"
		balanceFinal -= transaction.Value
	} else {
		eventType = "Credit"
		balanceFinal += transaction.Value
	}

	event := TransactionEvent{
		EventType:   eventType,
		Transaction: transaction,
	}

	s.ch <- event

	client := entity.Client{
		Limit:   balance.Limit,
		Balance: balanceFinal,
	}

	return client, nil
}

func (s *transactionService) FindBalance(ctx context.Context, clientID int) (entity.BankStatement, error) {
	if clientID < 1 || clientID > 5 {
		return entity.BankStatement{}, errors.New("no rows in result set")
	}

	return s.repo.FindBankStatement(ctx, clientID)
}
