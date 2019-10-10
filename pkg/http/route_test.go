package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouteServeHTTPNotFound(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost/hello", nil)
	router := Route()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, r)

	require.Equal(t, http.StatusNotFound, rec.Code)
}
