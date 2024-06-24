package services

import (
	"context"
	"feijuca/domain/entity"
	"feijuca/domain/ports/outbounds"
	"fmt"
)

type TransactionEvent struct {
	EventType   string
	Transaction entity.Transaction
}

type TransactionSubscriber struct {
	ch   chan TransactionEvent
	repo outbounds.Transaction
}

func NewTransactionSubscriber(repo outbounds.Transaction, ch chan TransactionEvent) *TransactionSubscriber {
	return &TransactionSubscriber{
		repo: repo,
		ch:   ch,
	}
}

func (s *TransactionSubscriber) Listen() {
	for {
		select {
		case event := <-s.ch:
			if event.EventType == "Debit" {
				_, err := s.repo.Save(context.Background(), event.Transaction)
				if err != nil {
					fmt.Println("ERROR SAVING DEBIT TRANSACTION")
					fmt.Println("error: ", err.Error())
				}
			}
			if event.EventType == "Credit" {
				_, err := s.repo.Save(context.Background(), event.Transaction)
				if err != nil {
					fmt.Println("ERROR SAVING CREDIT TRANSACTION")
					fmt.Println("error: ", err.Error())
				}
			}

		}
	}
}
