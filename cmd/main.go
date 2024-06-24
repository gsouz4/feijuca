package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"feijuca/api"
	"feijuca/domain/services"
	"feijuca/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/rinha?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
	)

	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch := make(chan services.TransactionEvent)

	repo := repository.NewTransactionRepository(conn)
	transactionService := services.NewTransactionService(repo, ch)
	transactionSubscriber := services.NewTransactionSubscriber(repo, ch)
	transactionController := api.NewTransactionController(transactionService)

	go transactionSubscriber.Listen()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())

	e.GET("/clientes/:id/extrato", transactionController.HandleFindStatement())
	e.POST("/clientes/:id/transacoes", transactionController.HandleCreateTransaction())

	e.Logger.Fatal(e.Start(":5000"))
}
