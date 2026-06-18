package server_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"goManip/jobs"
	"goManip/server"
)

type testStore struct{}

func (s testStore) SaveImage(ctx context.Context, jobId string, image []byte) error {
	return nil
}

func (s testStore) WriteError(ctx context.Context, jobId string, err error, traceId string) error {
	return nil
}

func (s testStore) GetImage(ctx context.Context, jobId string) ([]byte, string, error) {
	return nil, "", nil
}

func TestJobDispatcherMiddleware(t *testing.T) {
	e := echo.New()

	jd := &jobs.JobDispatcher{}
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
		wantErr      *server.GomanipError
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
				Trace:  "00000000000000000000000000000000",
			},
		},
		{
			name:         "missing content type",
			contentType:  "",
			expectStatus: http.StatusBadRequest,
			wantErr: &server.GomanipError{
				Status: "400",
				Detail: "Given filetype is not supported, please use jpeg or png images.",
				Trace:  "00000000000000000000000000000000",
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

func TestStoreMiddleware(t *testing.T) {
	e := echo.New()

	store := testStore{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	next := func(c echo.Context) error {
		val := c.Get("store")
		assert.NotNil(t, val, "store should be set in context")
		assert.Equal(t, store, val, "store should match the injected instance")
		return c.NoContent(http.StatusOK)
	}

	mw := server.StoreMiddleware(store)

	err := mw(next)(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
