package middleware


import (
	"fmt"
	"slices"
	"strings"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"goManip/JobDispatch"
	"goManip/errors"
)
var (
	supportedFileTypes = []string{ "image/png", "image/jpeg"}
)


func JobDispatcherMiddleware(jobDispatcher *JobDispatch.JobDispatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("jobDispatcher", jobDispatcher)
			return next(c)
		}
	}
}

func FileTypeVerifyMiddleware() echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc  {
		return func(c echo.Context) error {
			contentType:= c.Request().Header.Get("Content-Type")
			
			// get the actual filetype for logging and error messages
			// the Content-Type header has the format  type / subtype, however 
			// to be safe (in case its malformed), we get the file based on the 
			// last item in the split string.
			contents := strings.Split(contentType, "/")
			fileType := contents[len(contents)-1]

			if !slices.Contains(supportedFileTypes, contentType) {
				log.Error().Msg(fmt.Sprintf("request had content type of %s which is not supported", contentType))
				return errors.ReturnJsonError(c, http.StatusBadRequest, fmt.Sprintf("%s files are not supported", fileType))
			}

			return next(c)
		}
	}
}


