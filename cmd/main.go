package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"feijuca/domain/entity"
	"feijuca/repository"

	_ "github.com/lib/pq"
)

func main() {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/rinha?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewTransactionRepository(db)

	repo.Save(context.Background(), entity.Transaction{})
}
