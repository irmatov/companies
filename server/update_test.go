package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/irmatov/companies/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerUpdate(t *testing.T) {
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

	t.Run("get all companies", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		var r []types.Company
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, []types.Company{c1, c2}, r)
	})

	t.Run("bad url", func(t *testing.T) {
		c1.Name = "changed"
		b, err := json.Marshal(c1)
		require.NoError(t, err)

		req := httptest.NewRequest("PUT", baseURL+"garbage", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad content type", func(t *testing.T) {
		c1.Name = "changed"
		b, err := json.Marshal(c1)
		require.NoError(t, err)

		req := httptest.NewRequest("PUT", baseURL+strconv.Itoa(c1.Id), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/octet-stream")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("update first one", func(t *testing.T) {
		c1.Name = "changed"
		b, err := json.Marshal(c1)
		require.NoError(t, err)

		req := httptest.NewRequest("PUT", baseURL+strconv.Itoa(c1.Id), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("changes are visible", func(t *testing.T) {
		req := httptest.NewRequest("GET", baseURL, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		var r []types.Company
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, []types.Company{c1, c2}, r)
	})

	t.Run("update non-existing", func(t *testing.T) {
		c1.Id = 999
		b, err := json.Marshal(c1)
		require.NoError(t, err)

		req := httptest.NewRequest("PUT", baseURL+strconv.Itoa(c1.Id), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}
