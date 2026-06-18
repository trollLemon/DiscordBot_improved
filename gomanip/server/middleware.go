package server

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"

	"goManip/jobs"
	"goManip/worker"
)

var (
	supportedFileTypes = []string{"image/png", "image/jpeg"}
	ErrUnsupportedFile = errors.New("invalid filetype")
)

func JobDispatcherMiddleware(jobDispatcher *jobs.JobDispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("jobDispatcher", jobDispatcher)
			return next(c)
		}
	}
}

func StoreMiddleware(store worker.ImageStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("store", store)
			return next(c)
		}
	}
}

func FileTypeVerifyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")
			spanContext := trace.SpanFromContext(c.Request().Context()).SpanContext()
			if !slices.Contains(supportedFileTypes, contentType) {
				slog.ErrorContext(c.Request().Context(), fmt.Sprintf("request had content type of %s which is not supported", contentType))
				return SendGomanipError(c, ErrUnsupportedFile, spanContext.TraceID().String())
			}

			return next(c)
		}
	}
}
