package server

import (
	"errors"
	"fmt"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"goManip/jobdispatch"
)

var (
	supportedFileTypes = []string{"image/png", "image/jpeg"}
	ErrUnsupportedFile = errors.New("invalid filetype")
)

func JobDispatcherMiddleware(jobDispatcher *jobdispatch.JobDispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("jobDispatcher", jobDispatcher)
			return next(c)
		}
	}
}

func FileTypeVerifyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")

			if !slices.Contains(supportedFileTypes, contentType) {
				log.Error().Msg(fmt.Sprintf("request had content type of %s which is not supported", contentType))
				return SendGomanipError(c, ErrUnsupportedFile)
			}

			return next(c)
		}
	}
}
