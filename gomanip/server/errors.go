package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"goManip/jobs"
)

// GomanipError is the object returned to the API caller when an error occurs.
// It includes the status, and a detail string explaining the error.
type GomanipError struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}


// SendGomanipError returns a json payload containing the error message and status code.
func SendGomanipError(c echo.Context, err error ) error {
	errString := "an unexpected server side error occured."
	statusCode := http.StatusInternalServerError 

	var opErr *jobs.OperationParameterError
	
	if errors.As(err, &opErr) {
		errString = opErr.Error()	
		statusCode = http.StatusBadRequest 
	} else if errors.Is(err, jobs.ErrImgEmpty) {
		errString = "Given image was empty."
		statusCode = http.StatusUnprocessableEntity
	} else if errors.Is(err, jobs.ErrOpenCV) {
		errString = "OpenCV threw an error during image processing."
	} else if errors.Is(err, ErrUnsupportedFile) {
		errString = "Given filetype is not supported, please use jpeg or png images"
		statusCode = http.StatusBadRequest
	} 
	

	response := &GomanipError{
		Status: strconv.Itoa(statusCode),
		Detail: errString,
	}

	return c.JSON(statusCode, response)
}




