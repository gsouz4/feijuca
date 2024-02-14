package api

import (
	"context"
	"errors"
	"feijuca/domain/ports/inbounds"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type transactionController struct {
	service inbounds.TransactionService
}

type CreateRequest struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func (r *CreateRequest) Validate() error {
	if r.Type != "c" && r.Type != "d" {
		return errors.New("invalid request")
	}

	if len(r.Description) < 1 || len(r.Description) > 10 {
		return errors.New("invalid request")
	}

	return nil
}

func NewTransactionController(service inbounds.TransactionService) transactionController {
	return transactionController{
		service: service,
	}
}

func (c *transactionController) HandleCreateTransaction() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		context, cancel := context.WithCancel(context.Background())
		defer cancel()

		var input CreateRequest

		if err := ctx.Bind(&input); err != nil {
			return HandleError(ctx, errors.New("invalid request"))
		}

		if err := input.Validate(); err != nil {
			return HandleError(ctx, err)
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return HandleError(ctx, err)
		}

		client, err := c.service.Save(context, id, input.Value, input.Type, input.Description)
		if err != nil {
			log.Error(err)
			return HandleError(ctx, err)
		}

		return ctx.JSON(http.StatusOK, client)
	}
}

func (c *transactionController) HandleFindStatement() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		context, cancel := context.WithCancel(ctx.Request().Context())
		defer cancel()

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return HandleError(ctx, err)
		}

		statement, err := c.service.FindBalance(context, id)
		if err != nil {
			return HandleError(ctx, err)
		}

		return ctx.JSON(http.StatusOK, statement)
	}
}
