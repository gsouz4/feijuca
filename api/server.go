package api

import (
	"context"
	"errors"
	"feijuca/domain/ports/inbounds"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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
		context, cancel := context.WithCancel(ctx.Request().Context())
		defer cancel()

		var input CreateRequest

		if err := ctx.Bind(&input); err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		if err := input.Validate(); err != nil {
			return ctx.NoContent(http.StatusUnprocessableEntity)
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		if err := c.service.Save(context, id, input.Value, input.Type, input.Description); err != nil {
			return ctx.JSON(400, nil)
		}

		return ctx.NoContent(http.StatusOK)
	}
}

func (c *transactionController) HandleFindStatement() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		context, cancel := context.WithCancel(ctx.Request().Context())
		defer cancel()

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		statement, err := c.service.FindBalance(context, id)
		if err != nil {
			return ctx.NoContent(400) // TODO: fix
		}

		return ctx.JSON(http.StatusOK, statement)
	}
}
