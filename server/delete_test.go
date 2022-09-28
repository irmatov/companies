package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/irmatov/companies/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerDelete(t *testing.T) {
	db := getTestDatabase(t)
	server := New(db)
	c1 := types.Company{
		Name:    "First Company",
		Code:    "FIRST",
		Country: "UK",
		Website: "https://first.com/",
		Phone:   "+111",
	}
	c1.Id = testCreateCompany(t, server, c1)

	c2 := types.Company{
		Name:    "Second Company",
		Code:    "SECOND",
		Country: "FR",
		Website: "https://second.com/",
		Phone:   "+222",
	}
	c2.Id = testCreateCompany(t, server, c2)

	c3 := types.Company{
		Name:    "Third Company",
		Code:    "THIRD",
		Country: "PL",
		Website: "https://third.com/",
		Phone:   "+333",
	}
	c3.Id = testCreateCompany(t, server, c3)

	t.Run("get all companies", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		var r []types.Company
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, []types.Company{c1, c2, c3}, r)
	})

	t.Run("delete last one", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/companies/%d", c3.Id), nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("two companies present", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		var r []types.Company
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, []types.Company{c1, c2}, r)
	})

	t.Run("delete again, 404", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/companies/%d", c3.Id), nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
