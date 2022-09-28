package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/irmatov/companies/postgres"
	"github.com/irmatov/companies/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost/companies/"

func testCreateCompany(t *testing.T, h http.Handler, c types.Company) int {
	b, err := json.Marshal(c)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
	var r struct {
		Id int
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	return r.Id
}

func getTestDatabase(t *testing.T) types.Storage {
	db, err := sql.Open("pgx", "user=postgres password=postgres dbname=postgres host=localhost sslmode=disable")
	require.NoError(t, err)
	_, err = db.Exec("DROP TABLE IF EXISTS companies")
	require.NoError(t, err)
	_, err = db.Exec("DROP SEQUENCE IF EXISTS companies_id_seq")
	require.NoError(t, err)
	_, err = db.Exec("CREATE SEQUENCE companies_id_seq START 1")
	require.NoError(t, err)
	_, err = db.Exec(`
    CREATE TABLE companies (
        id INTEGER PRIMARY KEY DEFAULT nextval('companies_id_seq'),
        name TEXT NOT NULL UNIQUE,
        code TEXT NOT NULL,
        country TEXT NOT NULL,
        website TEXT,
        phone TEXT
    )`)
	require.NoError(t, err)
	return postgres.New(db)
}
