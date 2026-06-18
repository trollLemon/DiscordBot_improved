package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"goManip/jobs"
	"goManip/store"
)

// GomanipError is the object returned to the API caller when an error occurs.
type GomanipError struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
	Trace  string `json:"trace,omitempty"`
}

// SendGomanipError returns a json payload containing the error message and status code.
func SendGomanipError(c echo.Context, err error, trace string) error {
	errString := "an unexpected server side error occured."
	statusCode := http.StatusInternalServerError

	if errors.Is(err, ErrUnsupportedFile) {
		errString = "Given filetype is not supported, please use jpeg or png images."
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, jobs.ErrInvalidParameters) {
		errString = err.Error()
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, ErrParseParams) {
		errString = "Failed to parse query parameters. Check the request URI."
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, store.ErrImageNotFound) {
		errString = "No result found for the given job ID; it may still be processing or the ID is invalid."
		statusCode = http.StatusNotFound
	}

	response := &GomanipError{
		Status: strconv.Itoa(statusCode),
		Detail: errString,
		Trace:  trace,
	}

	return c.JSON(statusCode, response)
}
