package errors

import (
	"strconv"

	"github.com/labstack/echo/v4"
)


type GomanipError struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}


// JsonError returns a json payload containing the error message and status code.
func JsonError(c echo.Context, statusCode int, errString string) error {
	response := &GomanipError{
		Status: strconv.Itoa(statusCode),
		Detail: errString,
	}

	return c.JSON(statusCode, response)
}   
