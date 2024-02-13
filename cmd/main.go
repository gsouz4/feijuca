package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"feijuca/api"
	"feijuca/domain/services"
	"feijuca/repository"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	transactionService := services.NewTransactionRepository(repo)
	transactionController := api.NewTransactionController(transactionService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/clientes/:id/extrato", transactionController.HandleFindStatement())
	e.POST("/clientes/:id/transacoes", transactionController.HandleCreateTransaction())

	e.Logger.Fatal(e.Start(":8080"))
}
