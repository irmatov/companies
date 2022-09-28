package server

import (
	"net/http"
	"strconv"

	"github.com/irmatov/companies/types"
	"github.com/julienschmidt/httprouter"
)

func (s *server) delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeJson(w, http.StatusBadRequest, genericError{err.Error()})
		return
	}
	err = s.svc.Delete(r.Context(), id)
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
