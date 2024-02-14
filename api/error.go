package api

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	mappedErrors = map[string]int{
		"no rows in result set": http.StatusNotFound,
		"invalid transaction":   http.StatusUnprocessableEntity,
		"invalid request":       http.StatusUnprocessableEntity,
	}
)

func HandleError(c echo.Context, err error) error {
	code, ok := mappedErrors[err.Error()]
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(code)
}
