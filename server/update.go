package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/irmatov/companies/types"
	"github.com/julienschmidt/httprouter"
)

func (s *server) update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeJson(w, http.StatusBadRequest, genericError{err.Error()})
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeJson(w, http.StatusBadRequest, genericError{"invalid Content-Type"})
		return
	}
	var c types.Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeJson(w, http.StatusBadRequest, genericError{"invalid JSON"})
		return
	}
	if c.Id != id {
		writeJson(w, http.StatusBadRequest, genericError{"id mismatch"})
		return
	}
	err = s.svc.Update(r.Context(), c)
	if err != nil {
		if err == types.ErrNotFound {
			writeJson(w, http.StatusNotFound, genericError{err.Error()})
			return
		}
		writeJson(w, http.StatusInternalServerError, genericError{"internal server error"})
		return
	}
	writeJson(w, http.StatusNoContent, nil)
}
