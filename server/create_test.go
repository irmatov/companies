package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/irmatov/companies/filter"
	"github.com/irmatov/companies/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerCreate(t *testing.T) {
	// test creating a company
	db := getTestDatabase(t)
	server := New(db)
	c1 := types.Company{
		Name:    "test1",
		Code:    "TEST1",
		Country: "CN",
		Website: "http://test1.com",
		Phone:   "+111",
	}
	b, err := json.Marshal(c1)
	require.NoError(t, err)

	t.Run("wrong request content-type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/companies/", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/octet-stream")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
		var r genericError
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, "invalid Content-Type", r.Error)
	})

	t.Run("first create", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/companies/", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
		var r struct {
			Id int
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, 1, r.Id)
	})

	// attempt to create the same company again.
	// we accept the same company being created twice
	t.Run("second exact create request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/companies/", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
		var r struct {
			Id int
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
		assert.Equal(t, 1, r.Id)
	})

	// create another company with the same name but different code. should fail.
	t.Run("second exact create request", func(t *testing.T) {
		c2 := c1
		c2.Code = "TEST2"
		b, err := json.Marshal(c2)
		require.NoError(t, err)
		req := httptest.NewRequest("POST", "/companies/", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// finally, ensure that the company is present in the database
	t.Run("get company", func(t *testing.T) {
		_ = db.Tx(context.Background(), func(tx types.Tx) error {
			got, err := tx.Get(filter.Filter{})
			require.NoError(t, err)
			c1.Id = 1
			require.Equal(t, []types.Company{c1}, got)
			return nil
		})
	})
}
