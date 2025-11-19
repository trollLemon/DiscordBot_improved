package apierrors

import "errors"

var (
	ErrAPI           = errors.New("error returned from api")
	ErrServer        = errors.New("server side error")
	ErrNetwork       = errors.New("netork error")
	ErrResp          = errors.New("error marshaling response body")
	ErrRetry         = errors.New("retrying endpoint")
	ErrReading       = errors.New("error reading response body")
	ErrWriting       = errors.New("error writing multipart form data")
	ErrTimedOut      = errors.New("timed out calling service")
	ErrMakingRequest = errors.New("Error making http request")
)
