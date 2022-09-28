package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/irmatov/companies/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerGet(t *testing.T) {
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

	t.Run("get a single company", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL+strconv.Itoa(c1.Id), nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		var r types.Company
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, c1, r)
	})

	t.Run("get some companies with complex filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL, nil)
		q := req.URL.Query()
		q.Add("filter", `name,"First Company",=,phone,"+333",=,or`)
		req.URL.RawQuery = q.Encode()

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		var r []types.Company
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, []types.Company{c1, c3}, r)
	})

	t.Run("get non existing", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL+"999", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
