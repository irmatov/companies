package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type lookupResult struct {
	code        int
	countryCode string
}

func TestCountry(t *testing.T) {
	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	lookupResults := make(chan lookupResult, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := <-lookupResults
		w.WriteHeader(resp.code)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"country_code": "%s"}`, resp.countryCode)))
	}))
	defer server.Close()
	auth := &Country{
		Next:               protected,
		LookupURLFormat:    server.URL + "/%s",
		AllowedCountryCode: "US",
		Client:             &http.Client{},
	}

	t.Run("allowed", func(t *testing.T) {
		lookupResults <- lookupResult{http.StatusOK, "US"}
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "1.1.1.1:1234"
		w := httptest.NewRecorder()
		auth.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("denied", func(t *testing.T) {
		lookupResults <- lookupResult{http.StatusOK, "CN"}
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "2.2.2.2:1234"
		w := httptest.NewRecorder()
		auth.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, "forbidden\n", w.Body.String())
	})

	t.Run("lookup failure", func(t *testing.T) {
		lookupResults <- lookupResult{http.StatusInternalServerError, "CN"}
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "3.3.3.3:1234"
		w := httptest.NewRecorder()
		auth.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "internal server error\n", w.Body.String())
	})
}
