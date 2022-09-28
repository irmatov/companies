package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/irmatov/companies/types"
	"github.com/julienschmidt/httprouter"
)

// create will handle POST requests to /companies/
func (s *server) create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeJson(w, http.StatusBadRequest, genericError{"invalid Content-Type"})
		return
	}

	var c types.Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeJson(w, http.StatusBadRequest, genericError{"invalid JSON"})
		return
	}

	if c.Name == "" || strings.TrimSpace(c.Name) != c.Name {
		writeJson(w, http.StatusBadRequest, genericError{"company name is empty or contains leading/trailing spaces"})
		return
	}

	id, err := s.svc.Create(r.Context(), c)
	if err != nil {
		if err == types.ErrAlreadyExists {
			writeJson(w, http.StatusConflict, genericError{"company with the given name already exists"})
		} else {
			log.Printf("error creating company: %v", err)
			writeJson(w, http.StatusInternalServerError, genericError{"internal server error"})
		}
		return
	}
	writeJson(w, http.StatusCreated, struct {
		Id int
	}{id})
}
