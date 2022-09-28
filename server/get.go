package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/irmatov/companies/filter"
	"github.com/julienschmidt/httprouter"
)

// get will handle GET requests to /companies/
func (s *server) getMany(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var f filter.Filter
	var err error
	if expr := r.FormValue("filter"); expr != "" {
		f, err = filter.Execute(knownFields, expr)
		if err != nil {
			writeJson(w, http.StatusBadRequest, fmt.Sprintf("invalid filter expression: %s", err))
			return
		}
	}
	companies, err := s.svc.Get(r.Context(), f)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, genericError{"internal server error"})
		return
	}
	writeJson(w, http.StatusOK, companies)
}

func (s *server) getSingle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeJson(w, http.StatusBadRequest, genericError{err.Error()})
		return
	}
	f, err := filter.Execute(knownFields, fmt.Sprintf("id,%d,=", id))
	if err != nil {
		writeJson(w, http.StatusInternalServerError, genericError{"internal server error"})
	}
	result, err := s.svc.Get(r.Context(), f)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, genericError{"internal server error"})
		return
	}
	if len(result) == 0 {
		writeJson(w, http.StatusNotFound, genericError{"not found"})
		return
	}
	writeJson(w, http.StatusOK, result[0])
}
