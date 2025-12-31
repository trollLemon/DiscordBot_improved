package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"goManip/jobdispatch"
	"goManip/server"
)

func TestJobDispatcherMiddleware(t *testing.T) {
	e := echo.New()

	jd := &jobdispatch.JobDispatcher{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	next := func(c echo.Context) error {
		val := c.Get("jobDispatcher")
		assert.NotNil(t, val, "jobDispatcher should be set in context")
		assert.Equal(t, jd, val, "jobDispatcher should match the injected instance")
		return c.NoContent(http.StatusOK)
	}

	mw := server.JobDispatcherMiddleware(jd)

	err := mw(next)(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFileTypeVerifyMiddleware(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		contentType  string
		expectStatus int
		wantErr   *server.GomanipError
	}{
		{
			name:         "supported content type",
			contentType:  "image/png",
			expectStatus: http.StatusOK,
		},
		{
			name:         "unsupported content type",
			contentType:  "image/webp",
			expectStatus: http.StatusBadRequest,
			wantErr: &server.GomanipError{
				Status: "400",
				Detail: "Given filetype is not supported, please use jpeg or png images.",
			},
		},
		{
			name:         "missing content type",
			contentType:  "",
			expectStatus: http.StatusBadRequest,
			wantErr: &server.GomanipError{
				Status: "400",
				Detail: "Given filetype is not supported, please use jpeg or png images.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			next := func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}

			mw := server.FileTypeVerifyMiddleware()

			err := mw(next)(c)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectStatus, rec.Code)
			
			body, err := io.ReadAll(rec.Body)
			assert.Nil(t, err)
			
			if tt.wantErr != nil {
				var goManipErr *server.GomanipError
				err = json.Unmarshal(body, &goManipErr)
				assert.Nil(t, err)
				assert.Equal(t, tt.wantErr, goManipErr)
			}
		})
	}
}
